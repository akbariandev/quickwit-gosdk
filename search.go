package quickwitgosdk

type SearchRequest struct {
	Query          string   `json:"query"`
	SearchFields   []string `json:"search_field,omitempty"`
	StartTimestamp int64    `json:"start_timestamp,omitempty"`
	EndTimestamp   int64    `json:"end_timestamp,omitempty"`
	MaxHits        uint64   `json:"max_hits,omitempty"`
	StartOffset    uint64   `json:"start_offset,omitempty"`
	SortByField    string   `json:"sort_by_field,omitempty"`
}

type SearchResponse struct {
	NumHits           uint64        `json:"num_hits"`
	Hits              []interface{} `json:"hits"`
	ElapsedTimeMicros uint64        `json:"elapsed_time_micros"`
	Aggregations      interface{}   `json:"aggregations,omitempty"`
}

func (c *Client) Search(indexId string, searchRequest SearchRequest) (searchResponse SearchResponse, err error) {
	_, err = c.client.R().
		SetPathParam("indexId", indexId).
		SetBody(searchRequest).
		SetResult(&searchResponse).
		Post("/api/v1/{indexId}/search")

	return searchResponse, err
}
