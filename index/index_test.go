package index

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/akbariandev/quickwit-gosdk/client"
)

func TestList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{"version":"0.9","index_uid":"index-1:01ABC","index_config":{"version":"0.9","index_id":"index-1","index_uri":"file:///idx1"}},
			{"version":"0.9","index_uid":"index-2:01DEF","index_config":{"version":"0.9","index_id":"index-2","index_uri":"file:///idx2"}}
		]`))
	}))
	defer ts.Close()

	c := client.New(ts.URL)
	indexes, err := List(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(indexes) != 2 {
		t.Fatalf("expected 2 indexes, got %d", len(indexes))
	}
	if indexes[0].Config.IndexID != "index-1" {
		t.Errorf("expected index_id 'index-1', got %q", indexes[0].Config.IndexID)
	}
	if indexes[0].IndexUID != "index-1:01ABC" {
		t.Errorf("expected index_uid 'index-1:01ABC', got %q", indexes[0].IndexUID)
	}
	if indexes[1].Config.IndexURI != "file:///idx2" {
		t.Errorf("expected index_uri, got %q", indexes[1].Config.IndexURI)
	}
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/indexes/my-index" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"version":"0.9","index_uid":"my-index:01XYZ","index_config":{"version":"0.9","index_id":"my-index","index_uri":"file:///myidx"}}`))
	}))
	defer ts.Close()

	c := client.New(ts.URL)
	index, err := Get(c, "my-index")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if index.Config.IndexID != "my-index" {
		t.Errorf("expected index_id 'my-index', got %q", index.Config.IndexID)
	}
	if index.IndexUID != "my-index:01XYZ" {
		t.Errorf("expected index_uid 'my-index:01XYZ', got %q", index.IndexUID)
	}
}

func TestCreate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"version":"0.9","index_uid":"new-index:01NEW","index_config":{"version":"0.9","index_id":"new-index"}}`))
	}))
	defer ts.Close()

	c := client.New(ts.URL)
	index, err := Create(c, CreateRequest{IndexID: "new-index"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if index.Config.IndexID != "new-index" {
		t.Errorf("expected index_id 'new-index', got %q", index.Config.IndexID)
	}
}

func TestTimestampUnmarshalNumeric(t *testing.T) {
	// Quickwit often returns timestamps as numeric Unix seconds.
	var m Metadata
	if err := json.Unmarshal([]byte(`{"create_timestamp": 1704067200}`), &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Unix(1704067200, 0).UTC()
	if !m.CreateTimestamp.Time.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, m.CreateTimestamp.Time)
	}
}

func TestTimestampUnmarshalString(t *testing.T) {
	// Some Quickwit versions may return RFC3339 strings.
	var m Metadata
	if err := json.Unmarshal([]byte(`{"create_timestamp": "2024-01-01T00:00:00Z"}`), &m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !m.CreateTimestamp.Time.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, m.CreateTimestamp.Time)
	}
}

func TestDelete(t *testing.T) {
	// Quickwit returns an array for delete (empty for real, splits for dry-run).
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"split_id":"01ABC","num_docs":100}]`))
	}))
	defer ts.Close()

	c := client.New(ts.URL)
	resp, err := Delete(c, "old-index", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 split, got %d", len(resp))
	}
	if resp[0].SplitID != "01ABC" {
		t.Errorf("expected split_id '01ABC', got %q", resp[0].SplitID)
	}
}

func TestDeleteEmpty(t *testing.T) {
	// Real delete returns an empty array.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	c := client.New(ts.URL)
	resp, err := Delete(c, "old-index", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp) != 0 {
		t.Errorf("expected 0 splits, got %d", len(resp))
	}
}
