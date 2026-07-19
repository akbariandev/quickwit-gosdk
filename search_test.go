package quickwitgosdk

import (
	"context"
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{
			NumHits:           2,
			ElapsedTimeMicros: 1500,
			Hits: []Hit{
				{JSON: map[string]interface{}{"title": "doc1"}},
				{JSON: map[string]interface{}{"title": "doc2"}},
			},
		})
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

func TestSearchStreamCallback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")

		r1 := SearchResponse{NumHits: 1, Hits: []Hit{{JSON: map[string]string{"id": "1"}}}}
		r2 := SearchResponse{NumHits: 2, Hits: []Hit{{JSON: map[string]string{"id": "2"}}}}

		json.NewEncoder(w).Encode(r1)
		json.NewEncoder(w).Encode(r2)
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	var results []SearchResponse

	err := client.SearchStreamWithCallback(context.Background(), "my-index", SearchStreamRequest{
		SearchRequest: SearchRequest{Query: "test"},
	}, func(resp SearchResponse) error {
		results = append(results, resp)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].NumHits != 1 || results[1].NumHits != 2 {
		t.Error("unexpected num_hits in stream responses")
	}
}

func TestSearchScrollChannel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api/v1/my-index/search/scroll") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/x-ndjson")

		r1 := SearchResponse{NumHits: 1, ElapsedTimeMicros: 100}
		r2 := SearchResponse{NumHits: 2, ElapsedTimeMicros: 200}

		json.NewEncoder(w).Encode(r1)
		json.NewEncoder(w).Encode(r2)
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	respCh, errCh := client.SearchScroll(context.Background(), "my-index", SearchScrollRequest{
		SearchRequest: SearchRequest{Query: "test"},
		ScrollTTLSecs: 60,
	})

	var results []SearchResponse
	for resp := range respCh {
		results = append(results, resp)
	}

	if err := <-errCh; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
