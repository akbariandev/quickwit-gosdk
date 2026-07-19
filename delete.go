package quickwitgosdk

import "time"

// DeleteQueryRequest is the request body for submitting a delete-by-query task.
type DeleteQueryRequest struct {
	Query          string   `json:"query"`
	SearchFields   []string `json:"search_field,omitempty"`
	StartTimestamp *int64   `json:"start_timestamp,omitempty"`
	EndTimestamp   *int64   `json:"end_timestamp,omitempty"`
	TagFilters     []string `json:"tag_filters,omitempty"`
	Filter         string   `json:"filter,omitempty"`
}

// DeleteQueryResponse is the response returned after submitting a delete-by-query task.
type DeleteQueryResponse struct {
	TaskID string `json:"task_id"`
}

// DeleteTaskResponse is the response returned when querying the status of a delete task.
type DeleteTaskResponse struct {
	TaskID    string     `json:"task_id"`
	Status    string     `json:"status"`     // "running", "success", "error", "cancelled"
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Error     string     `json:"error,omitempty"`
}

// DeleteByQuery submits a delete-by-query task for the given index.
func (c *Client) DeleteByQuery(indexId string, req DeleteQueryRequest) (DeleteQueryResponse, error) {
	var resp DeleteQueryResponse
	_, err := c.client.R().
		SetPathParam("indexId", indexId).
		SetBody(req).
		SetResult(&resp).
		Post("/api/v1/{indexId}/delete-tasks")

	return resp, err
}

// GetDeleteTask returns the status of a delete task.
func (c *Client) GetDeleteTask(indexId string, taskId string) (DeleteTaskResponse, error) {
	var resp DeleteTaskResponse
	_, err := c.client.R().
		SetPathParam("indexId", indexId).
		SetPathParam("taskId", taskId).
		SetResult(&resp).
		Get("/api/v1/{indexId}/delete-tasks/{taskId}")

	return resp, err
}
