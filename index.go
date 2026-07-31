package quickwitgosdk

import (
	"encoding/json"
	"time"
)

// Timestamp is a custom type that can unmarshal from both numeric (Unix seconds)
// and RFC3339 string timestamps returned by the Quickwit API.
type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Try numeric first (e.g. 1704067200)
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		secs := int64(num)
		nanos := int64((num - float64(secs)) * 1e9)
		t.Time = time.Unix(secs, nanos).UTC()
		return nil
	}

	// Fall back to RFC3339 string
	return json.Unmarshal(data, &t.Time)
}

// IndexMetadata represents the top-level metadata for a Quickwit index.
type IndexMetadata struct {
	Version         string                          `json:"version"`
	IndexUID        string                          `json:"index_uid"`
	IndexConfig     IndexConfig                     `json:"index_config"`
	Checkpoint      map[string]map[string]interface{} `json:"checkpoint,omitempty"`
	CreateTimestamp Timestamp                       `json:"create_timestamp,omitempty"`
	Sources         []Source                        `json:"sources,omitempty"`
}

// IndexConfig represents the configuration of a Quickwit index.
type IndexConfig struct {
	Version          string           `json:"version"`
	IndexID          string           `json:"index_id"`
	IndexURI         string           `json:"index_uri,omitempty"`
	DocMapping       DocMapping       `json:"doc_mapping,omitempty"`
	IndexingSettings IndexingSettings `json:"indexing_settings,omitempty"`
	IngestSettings   IngestSettings   `json:"ingest_settings,omitempty"`
	SearchSettings   SearchSettings   `json:"search_settings,omitempty"`
	Retention        *Retention       `json:"retention,omitempty"`
}

// CreateIndexRequest is the request body for creating a new index.
type CreateIndexRequest struct {
	Version          string           `json:"version,omitempty"`
	IndexID          string           `json:"index_id"`
	IndexURI         string           `json:"index_uri,omitempty"`
	DocMapping       DocMapping       `json:"doc_mapping,omitempty"`
	IndexingSettings IndexingSettings `json:"indexing_settings,omitempty"`
	IngestSettings   IngestSettings   `json:"ingest_settings,omitempty"`
	SearchSettings   SearchSettings   `json:"search_settings,omitempty"`
	Retention        *Retention       `json:"retention,omitempty"`
	Overrides        json.RawMessage  `json:"overrides,omitempty"`
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
