package quickwitgosdk

// SearchRequest is the request body for the Quickwit search API.
type SearchRequest struct {
	Query               string                 `json:"query"`
	DefaultOperator     string                 `json:"default_operator,omitempty"`     // "AND" or "OR"
	SearchFields        []string               `json:"search_field,omitempty"`
	StartTimestamp      *int64                 `json:"start_timestamp,omitempty"`       // use pointer to distinguish zero value
	EndTimestamp        *int64                 `json:"end_timestamp,omitempty"`
	MaxHits             uint64                 `json:"max_hits,omitempty"`
	StartOffset         uint64                 `json:"start_offset,omitempty"`
	SortByField         *SortByField           `json:"sort_by_field,omitempty"`
	SortByFieldDocOrder *SortByFieldDocOrder   `json:"sort_by_field_doc_order,omitempty"`
	Aggregations        map[string]interface{} `json:"aggregations,omitempty"`
	Source              string                 `json:"_source,omitempty"`
	SnippetFields       *SnippetRequest        `json:"snippet_fields,omitempty"`
	HighlightFields     *HighlightRequest       `json:"highlight_fields,omitempty"`
	TagFilters          []string               `json:"tag_filters,omitempty"`
	Filter              string                 `json:"filter,omitempty"`
}

// SearchResponse is the response body for the Quickwit search API.
type SearchResponse struct {
	NumHits           uint64        `json:"num_hits"`
	Hits              []Hit         `json:"hits"`
	ElapsedTimeMicros uint64        `json:"elapsed_time_micros"`
	Errors            []SearchError `json:"errors,omitempty"`
	Aggregations      interface{}   `json:"aggregations,omitempty"`
}

// Search performs a search request against the given index.
func (c *Client) Search(indexId string, searchRequest SearchRequest) (searchResponse SearchResponse, err error) {
	_, err = c.client.R().
		SetPathParam("indexId", indexId).
		SetBody(searchRequest).
		SetResult(&searchResponse).
		Post("/api/v1/{indexId}/search")

	return searchResponse, err
}
