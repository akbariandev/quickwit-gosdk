// Package index provides the Quickwit index management API.
package index

import (
	"encoding/json"

	"github.com/akbariandev/quickwit-gosdk/client"
	"github.com/akbariandev/quickwit-gosdk/internal/httputil"
	"github.com/akbariandev/quickwit-gosdk/types"
)

// Metadata represents the top-level metadata for a Quickwit index.
type Metadata struct {
	Version         string                            `json:"version"`
	IndexUID        string                            `json:"index_uid"`
	Config          Config                            `json:"index_config"`
	Checkpoint      map[string]map[string]interface{} `json:"checkpoint,omitempty"`
	CreateTimestamp types.Timestamp                   `json:"create_timestamp,omitempty"`
	Sources         []types.Source                    `json:"sources,omitempty"`
}

// Config represents the configuration of a Quickwit index.
type Config struct {
	Version          string                 `json:"version"`
	IndexID          string                 `json:"index_id"`
	IndexURI         string                 `json:"index_uri,omitempty"`
	DocMapping       types.DocMapping       `json:"doc_mapping,omitempty"`
	IndexingSettings types.IndexingSettings `json:"indexing_settings,omitempty"`
	IngestSettings   types.IngestSettings   `json:"ingest_settings,omitempty"`
	SearchSettings   types.SearchSettings   `json:"search_settings,omitempty"`
	Retention        *types.Retention       `json:"retention,omitempty"`
}

// CreateRequest is the request body for creating a new index.
// Version, IndexID, and DocMapping are required by the Quickwit API.
type CreateRequest struct {
	Version          string                 `json:"version,omitempty"`
	IndexID          string                 `json:"index_id"`
	IndexURI         string                 `json:"index_uri,omitempty"`
	DocMapping       types.DocMapping       `json:"doc_mapping,omitempty"`
	IndexingSettings types.IndexingSettings `json:"indexing_settings,omitempty"`
	IngestSettings   types.IngestSettings   `json:"ingest_settings,omitempty"`
	SearchSettings   types.SearchSettings   `json:"search_settings,omitempty"`
	Retention        *types.Retention       `json:"retention,omitempty"`
	Overrides        json.RawMessage        `json:"overrides,omitempty"`
}

// DeleteResponse represents a single split entry returned when deleting an index.
type DeleteResponse struct {
	SplitID                   string `json:"split_id"`
	NumDocs                   int    `json:"num_docs"`
	UncompressedDocsSizeBytes string `json:"uncompressed_docs_size_bytes,omitempty"`
	FileName                  string `json:"file_name,omitempty"`
	FileSizeBytes             string `json:"file_size_bytes,omitempty"`
}

// Create creates a new Quickwit index.
func Create(c *client.Client, req CreateRequest) (Metadata, error) {
	var resp Metadata
	_, err := c.HTTP.R().
		SetBody(req).
		SetResult(&resp).
		Post("/api/v1/indexes")

	return resp, err
}

// Get returns the metadata for a specific index.
func Get(c *client.Client, indexId string) (Metadata, error) {
	var resp Metadata
	_, err := c.HTTP.R().
		SetPathParam("indexId", indexId).
		SetResult(&resp).
		Get("/api/v1/indexes/{indexId}")

	return resp, err
}

// List returns metadata for all indexes.
func List(c *client.Client) ([]Metadata, error) {
	var resp []Metadata
	_, err := c.HTTP.R().
		SetResult(&resp).
		Get("/api/v1/indexes")

	return resp, err
}

// Delete deletes an index. If dryRun is true, the deletion is only simulated.
// Quickwit returns an array of affected splits (empty for a real delete).
func Delete(c *client.Client, indexId string, dryRun bool) ([]DeleteResponse, error) {
	var resp []DeleteResponse
	_, err := c.HTTP.R().
		SetPathParam("indexId", indexId).
		SetQueryParam("dry_run", httputil.BoolToString(dryRun)).
		SetResult(&resp).
		Delete("/api/v1/indexes/{indexId}")

	return resp, err
}

// Clear removes all splits from an index without deleting the index itself.
func Clear(c *client.Client, indexId string) error {
	_, err := c.HTTP.R().
		SetPathParam("indexId", indexId).
		Post("/api/v1/indexes/{indexId}/clear")

	return err
}
