// Package search provides the Quickwit search API.
package search

import (
	"github.com/akbariandev/quickwit-gosdk/client"
	"github.com/akbariandev/quickwit-gosdk/types"
)

// Request is the request body for the Quickwit search API.
type Request struct {
	Query               string                     `json:"query"`
	DefaultOperator     string                     `json:"default_operator,omitempty"` // "AND" or "OR"
	SearchFields        []string                   `json:"search_field,omitempty"`
	StartTimestamp      *int64                     `json:"start_timestamp,omitempty"` // use pointer to distinguish zero value
	EndTimestamp        *int64                     `json:"end_timestamp,omitempty"`
	MaxHits             uint64                     `json:"max_hits,omitempty"`
	StartOffset         uint64                     `json:"start_offset,omitempty"`
	SortByField         *types.SortByField         `json:"sort_by_field,omitempty"`
	SortByFieldDocOrder *types.SortByFieldDocOrder `json:"sort_by_field_doc_order,omitempty"`
	Aggregations        map[string]interface{}     `json:"aggs,omitempty"`
	Source              string                     `json:"_source,omitempty"`
	SnippetFields       *types.SnippetRequest      `json:"snippet_fields,omitempty"`
	HighlightFields     *types.HighlightRequest    `json:"highlight_fields,omitempty"`
	TagFilters          []string                   `json:"tag_filters,omitempty"`
	Filter              string                     `json:"filter,omitempty"`
	Format              types.OutputFormat         `json:"format,omitempty"`          // "json" or "pretty_json"
	ScrollTTLSecs       *uint64                    `json:"scroll_ttl_secs,omitempty"` // enables scroll mode when set
}

// Response is the response body for the Quickwit search API.
type Response struct {
	NumHits           uint64              `json:"num_hits"`
	Hits              []types.Hit         `json:"hits"`
	ElapsedTimeMicros uint64              `json:"elapsed_time_micros"`
	Errors            []types.SearchError `json:"errors,omitempty"`
	Aggregations      interface{}         `json:"aggregations,omitempty"`
}

// Do performs a search request against the given index.
func Do(c *client.Client, indexId string, req Request) (Response, error) {
	var resp Response
	_, err := c.HTTP.R().
		SetPathParam("indexId", indexId).
		SetBody(req).
		SetResult(&resp).
		Post("/api/v1/{indexId}/search")

	return resp, err
}
