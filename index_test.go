package quickwitgosdk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListIndexes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]IndexMetadata{
			{IndexID: "index-1"},
			{IndexID: "index-2"},
		})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	indexes, err := client.ListIndexes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(indexes) != 2 {
		t.Errorf("expected 2 indexes, got %d", len(indexes))
	}
	if indexes[0].IndexID != "index-1" {
		t.Errorf("expected index_id 'index-1', got %q", indexes[0].IndexID)
	}
}

func TestGetIndex(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GET /api/v1/indexes/{indexId}
		if r.URL.Path != "/api/v1/indexes/my-index" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(IndexMetadata{IndexID: "my-index"})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	index, err := client.GetIndex("my-index")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if index.IndexID != "my-index" {
		t.Errorf("expected index_id 'my-index', got %q", index.IndexID)
	}
}

func TestCreateIndex(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(IndexMetadata{IndexID: "new-index"})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	index, err := client.CreateIndex(CreateIndexRequest{IndexID: "new-index"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if index.IndexID != "new-index" {
		t.Errorf("expected index_id 'new-index', got %q", index.IndexID)
	}
}

func TestTimestampUnmarshalNumeric(t *testing.T) {
	// Quickwit often returns timestamps as numeric Unix seconds.
	data := []byte(`{"create_timestamp": 1704067200}`)
	var m IndexMetadata
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Unix(1704067200, 0).UTC()
	if !m.CreateTimestamp.Time.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, m.CreateTimestamp.Time)
	}
}

func TestTimestampUnmarshalString(t *testing.T) {
	// Some Quickwit versions may return RFC3339 strings.
	data := []byte(`{"create_timestamp": "2024-01-01T00:00:00Z"}`)
	var m IndexMetadata
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !m.CreateTimestamp.Time.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, m.CreateTimestamp.Time)
	}
}

func TestDeleteIndex(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeleteIndexResponse{IndexID: "old-index"})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	resp, err := client.DeleteIndex("old-index", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.IndexID != "old-index" {
		t.Errorf("expected index_id 'old-index', got %q", resp.IndexID)
	}
}
