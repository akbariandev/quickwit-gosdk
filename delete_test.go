package quickwitgosdk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteByQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/api/v1/my-index/delete-tasks") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeleteQueryResponse{TaskID: "task-123"})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	resp, err := client.DeleteByQuery("my-index", DeleteQueryRequest{Query: "old logs"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TaskID != "task-123" {
		t.Errorf("expected task_id 'task-123', got %q", resp.TaskID)
	}
}

func TestGetDeleteTask(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/api/v1/my-index/delete-tasks/task-123") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeleteTaskResponse{
			TaskID: "task-123",
			Status: "success",
		})
	}))
	defer ts.Close()

	client := NewClient(ts.URL)
	resp, err := client.GetDeleteTask("my-index", "task-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected status 'success', got %q", resp.Status)
	}
}
