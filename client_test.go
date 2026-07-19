package quickwitgosdk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClientWithOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-api-key" {
			t.Errorf("expected Authorization header 'Bearer test-api-key', got %q", auth)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{NumHits: 0})
	}))
	defer ts.Close()

	client := NewClient(ts.URL,
		WithAPIKey("test-api-key"),
		WithTimeout(30*time.Second),
	)

	resp, err := client.Search("my-index", SearchRequest{Query: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.NumHits != 0 {
		t.Errorf("expected 0 hits, got %d", resp.NumHits)
	}
}

func TestNewClientWithTransport(t *testing.T) {
	var gotRequest bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequest = true
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{NumHits: 1})
	}))
	defer ts.Close()

	client := NewClient(ts.URL, WithTransport(ts.Client().Transport))

	_, err := client.Search("idx", SearchRequest{Query: "q"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !gotRequest {
		t.Error("expected request to reach server")
	}
}

func TestQuickwitError(t *testing.T) {
	e := &QuickwitError{StatusCode: 404, Message: "index not found"}
	if e.Error() != "quickwit api error (status 404): index not found" {
		t.Errorf("unexpected error message: %s", e.Error())
	}
}
