package delete

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akbariandev/quickwit-gosdk/client"
)

func TestSubmit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/v1/my-index/delete-tasks") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"task_id":"task-123"}`))
	}))
	defer ts.Close()

	c := client.New(ts.URL)
	resp, err := Submit(c, "my-index", Request{Query: "old logs"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TaskID != "task-123" {
		t.Errorf("expected task_id 'task-123', got %q", resp.TaskID)
	}
}

func TestGetTask(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api/v1/my-index/delete-tasks/task-123") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"task_id":"task-123","status":"success"}`))
	}))
	defer ts.Close()

	c := client.New(ts.URL)
	resp, err := GetTask(c, "my-index", "task-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected status 'success', got %q", resp.Status)
	}
}
