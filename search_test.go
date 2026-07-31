package quickwitgosdk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/v1/my-index/search") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Query != "hello world" {
			t.Errorf("expected query 'hello world', got %q", req.Query)
		}
		if req.DefaultOperator != "AND" {
			t.Errorf("expected default_operator 'AND', got %q", req.DefaultOperator)
		}

		// Return flat document fields as Quickwit does
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"num_hits": 2,
			"elapsed_time_micros": 1500,
			"hits": [
				{"title": "doc1", "body": "first doc"},
				{"title": "doc2", "body": "second doc"}
			]
		}`))
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	resp, err := client.Search("my-index", SearchRequest{
		Query:           "hello world",
		DefaultOperator: "AND",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.NumHits != 2 {
		t.Errorf("expected 2 hits, got %d", resp.NumHits)
	}
	if len(resp.Hits) != 2 {
		t.Errorf("expected 2 hits, got %d", len(resp.Hits))
	}
	if resp.Hits[0].Fields["title"] != "doc1" {
		t.Errorf("expected title 'doc1', got %v", resp.Hits[0].Fields["title"])
	}
	if resp.Hits[1].Fields["body"] != "second doc" {
		t.Errorf("expected body 'second doc', got %v", resp.Hits[1].Fields["body"])
	}
}

func TestSearchWithSortAndSnippets(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req SearchRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.SortByField == nil || req.SortByField.FieldName != "timestamp" {
			t.Error("expected sort_by_field with field_name 'timestamp'")
		}
		if req.SnippetFields == nil || req.SnippetFields.FieldName != "body" {
			t.Error("expected snippet_fields with field_name 'body'")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{NumHits: 0})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	_, err := client.Search("idx", SearchRequest{
		Query:       "test",
		SortByField: &SortByField{FieldName: "timestamp", Order: "desc"},
		SnippetFields: &SnippetRequest{
			FieldName:              "body",
			MaxNumCharsPerFragment: 200,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSearchWithScroll(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/v1/my-index/search") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.ScrollTTLSecs == nil || *req.ScrollTTLSecs != 60 {
			t.Errorf("expected scroll_ttl_secs=60, got %v", req.ScrollTTLSecs)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{NumHits: 0})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	scrollTTL := uint64(60)
	resp, err := client.Search("my-index", SearchRequest{
		Query:         "test",
		ScrollTTLSecs: &scrollTTL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.NumHits != 0 {
		t.Errorf("expected 0 hits, got %d", resp.NumHits)
	}
}
