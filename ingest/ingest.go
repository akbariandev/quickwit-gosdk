// Package ingest provides the Quickwit document ingestion API.
package ingest

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/akbariandev/quickwit-gosdk/client"
)

// Response is the response returned after ingesting documents.
type Response struct {
	NumPersisted uint64       `json:"num_persisted"`
	NumFailed    uint64       `json:"num_failed,omitempty"`
	Errors       []BatchError `json:"errors,omitempty"`
}

// BatchError represents an error for a single document in a batch.
type BatchError struct {
	DocJSON interface{} `json:"doc_json,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Send sends a batch of documents to the given index as NDJSON.
func Send(c *client.Client, indexId string, docs []interface{}) (Response, error) {
	var resp Response

	ndjson, err := marshalNDJSON(docs)
	if err != nil {
		return resp, err
	}

	_, err = c.HTTP.R().
		SetPathParam("indexId", indexId).
		SetHeader("Content-Type", "application/x-ndjson").
		SetBody(ndjson).
		SetResult(&resp).
		Post("/api/v1/ingest/{indexId}")

	return resp, err
}

// SendFromReader sends documents from an io.Reader to the given index.
// The reader should provide NDJSON-formatted data.
func SendFromReader(c *client.Client, indexId string, reader io.Reader) (Response, error) {
	var resp Response

	_, err := c.HTTP.R().
		SetPathParam("indexId", indexId).
		SetHeader("Content-Type", "application/x-ndjson").
		SetBody(reader).
		SetResult(&resp).
		Post("/api/v1/ingest/{indexId}")

	return resp, err
}

// ForceMerge triggers a force-merge operation on the given index.
func ForceMerge(c *client.Client, indexId string) error {
	_, err := c.HTTP.R().
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
