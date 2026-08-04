// Package delete provides the Quickwit delete-by-query task API.
package delete

import (
	"time"

	"github.com/akbariandev/quickwit-gosdk/client"
)

// Request is the request body for submitting a delete-by-query task.
type Request struct {
	Query          string   `json:"query"`
	SearchFields   []string `json:"search_field,omitempty"`
	StartTimestamp *int64   `json:"start_timestamp,omitempty"`
	EndTimestamp   *int64   `json:"end_timestamp,omitempty"`
	TagFilters     []string `json:"tag_filters,omitempty"`
	Filter         string   `json:"filter,omitempty"`
}

// Response is the response returned after submitting a delete-by-query task.
type Response struct {
	TaskID string `json:"task_id"`
}

// TaskResponse is the response returned when querying the status of a delete task.
type TaskResponse struct {
	TaskID    string    `json:"task_id"`
	Status    string    `json:"status"` // "running", "success", "error", "cancelled"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Error     string    `json:"error,omitempty"`
}

// Submit submits a delete-by-query task for the given index.
func Submit(c *client.Client, indexId string, req Request) (Response, error) {
	var resp Response
	_, err := c.HTTP.R().
		SetPathParam("indexId", indexId).
		SetBody(req).
		SetResult(&resp).
		Post("/api/v1/{indexId}/delete-tasks")

	return resp, err
}

// GetTask returns the status of a delete task.
func GetTask(c *client.Client, indexId string, taskId string) (TaskResponse, error) {
	var resp TaskResponse
	_, err := c.HTTP.R().
		SetPathParam("indexId", indexId).
		SetPathParam("taskId", taskId).
		SetResult(&resp).
		Get("/api/v1/{indexId}/delete-tasks/{taskId}")

	return resp, err
}
