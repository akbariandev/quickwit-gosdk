package quickwitgosdk

import (
	"bytes"
	"encoding/json"
	"io"
)

// IngestResponse is the response returned after ingesting documents.
type IngestResponse struct {
	NumPersisted uint64              `json:"num_persisted"`
	NumFailed    uint64              `json:"num_failed,omitempty"`
	Errors       []IngestBatchError  `json:"errors,omitempty"`
}

// IngestBatchError represents an error for a single document in a batch.
type IngestBatchError struct {
	DocJSON interface{} `json:"doc_json,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Ingest sends a batch of documents to the given index as NDJSON.
func (c *Client) Ingest(indexId string, docs []interface{}) (IngestResponse, error) {
	var resp IngestResponse

	ndjson, err := marshalNDJSON(docs)
	if err != nil {
		return resp, err
	}

	_, err = c.client.R().
		SetPathParam("indexId", indexId).
		SetHeader("Content-Type", "application/x-ndjson").
		SetBody(ndjson).
		SetResult(&resp).
		Post("/api/v1/ingest/{indexId}")

	return resp, err
}

// IngestFromReader sends documents from an io.Reader to the given index.
// The reader should provide NDJSON-formatted data.
func (c *Client) IngestFromReader(indexId string, reader io.Reader) (IngestResponse, error) {
	var resp IngestResponse

	_, err := c.client.R().
		SetPathParam("indexId", indexId).
		SetHeader("Content-Type", "application/x-ndjson").
		SetBody(reader).
		SetResult(&resp).
		Post("/api/v1/ingest/{indexId}")

	return resp, err
}

// ForceMerge triggers a force-merge operation on the given index.
func (c *Client) ForceMerge(indexId string) error {
	_, err := c.client.R().
		SetPathParam("indexId", indexId).
		Post("/api/v1/indexes/{indexId}/force-merge")

	return err
}

// marshalNDJSON encodes a slice of documents as newline-delimited JSON.
func marshalNDJSON(docs []interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, doc := range docs {
		if err := enc.Encode(doc); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
