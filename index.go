package quickwitgosdk

import (
	"encoding/json"
	"time"
)

// IndexMetadata represents metadata about a Quickwit index.
type IndexMetadata struct {
	IndexID         string            `json:"index_id"`
	URI             string            `json:"index_uri,omitempty"`
	Version         string            `json:"version,omitempty"`
	DocMapping      json.RawMessage   `json:"doc_mapping,omitempty"`
	IndexingSettings json.RawMessage  `json:"indexing_settings,omitempty"`
	SearchSettings  json.RawMessage   `json:"search_settings,omitempty"`
	Retention       json.RawMessage   `json:"retention,omitempty"`
	CreateTimestamp  time.Time         `json:"create_timestamp,omitempty"`
	Source          json.RawMessage   `json:"source,omitempty"`
}

// CreateIndexRequest is the request body for creating a new index.
type CreateIndexRequest struct {
	IndexID         string          `json:"index_id"`
	DocMapping      json.RawMessage `json:"doc_mapping,omitempty"`
	IndexingSettings json.RawMessage `json:"indexing_settings,omitempty"`
	SearchSettings  json.RawMessage `json:"search_settings,omitempty"`
	Retention       json.RawMessage `json:"retention,omitempty"`
	Source          json.RawMessage `json:"source,omitempty"`
	Overrides       json.RawMessage `json:"overrides,omitempty"`
}

// DeleteIndexResponse is the response returned after deleting an index.
type DeleteIndexResponse struct {
	IndexID string `json:"index_id"`
}

// CreateIndex creates a new Quickwit index.
func (c *Client) CreateIndex(req CreateIndexRequest) (IndexMetadata, error) {
	var resp IndexMetadata
	_, err := c.client.R().
		SetBody(req).
		SetResult(&resp).
		Post("/api/v1/indexes")

	return resp, err
}

// GetIndex returns the metadata for a specific index.
func (c *Client) GetIndex(indexId string) (IndexMetadata, error) {
	var resp IndexMetadata
	_, err := c.client.R().
		SetPathParam("indexId", indexId).
		SetResult(&resp).
		Get("/api/v1/indexes/{indexId}")

	return resp, err
}

// ListIndexes returns metadata for all indexes.
func (c *Client) ListIndexes() ([]IndexMetadata, error) {
	var resp []IndexMetadata
	_, err := c.client.R().
		SetResult(&resp).
		Get("/api/v1/indexes")

	return resp, err
}

// DeleteIndex deletes an index. If dryRun is true, the deletion is only simulated.
func (c *Client) DeleteIndex(indexId string, dryRun bool) (DeleteIndexResponse, error) {
	var resp DeleteIndexResponse
	_, err := c.client.R().
		SetPathParam("indexId", indexId).
		SetQueryParam("dry_run", boolToString(dryRun)).
		SetResult(&resp).
		Delete("/api/v1/indexes/{indexId}")

	return resp, err
}

// ClearIndex removes all splits from an index without deleting the index itself.
func (c *Client) ClearIndex(indexId string) error {
	_, err := c.client.R().
		SetPathParam("indexId", indexId).
		Post("/api/v1/indexes/{indexId}/clear")

	return err
}

// boolToString converts a bool to its string representation for query parameters.
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
