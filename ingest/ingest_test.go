package ingest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akbariandev/quickwit-gosdk/client"
)

func TestSend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/ingest/my-index" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/x-ndjson" {
			t.Errorf("expected Content-Type application/x-ndjson, got %q", ct)
		}

		// Verify body is NDJSON
		decoder := json.NewDecoder(r.Body)
		count := 0
		for decoder.More() {
			var doc map[string]interface{}
			if err := decoder.Decode(&doc); err != nil {
				t.Fatalf("failed to decode doc: %v", err)
			}
			count++
		}
		if count != 2 {
			t.Errorf("expected 2 docs in body, got %d", count)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"num_persisted":2}`))
	}))
	defer ts.Close()

	c := client.New(ts.URL)
	resp, err := Send(c, "my-index", []interface{}{
		map[string]string{"title": "doc1"},
		map[string]string{"title": "doc2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.NumPersisted != 2 {
		t.Errorf("expected 2 persisted, got %d", resp.NumPersisted)
	}
}

func TestSendFromReader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"num_persisted":1}`))
	}))
	defer ts.Close()

	c := client.New(ts.URL)
	ndjson := `{"title":"doc1"}
{"title":"doc2"}
`
	resp, err := SendFromReader(c, "my-index", strings.NewReader(ndjson))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.NumPersisted != 1 {
		t.Errorf("expected 1 persisted, got %d", resp.NumPersisted)
	}
}

func TestMarshalNDJSON(t *testing.T) {
	docs := []interface{}{
		map[string]string{"a": "1"},
		map[string]int{"b": 2},
	}

	data, err := marshalNDJSON(docs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	var first map[string]string
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("failed to decode first line: %v", err)
	}
	if first["a"] != "1" {
		t.Errorf("expected first doc a=1, got %q", first["a"])
	}
}
