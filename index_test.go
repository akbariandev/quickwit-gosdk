package quickwitgosdk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
