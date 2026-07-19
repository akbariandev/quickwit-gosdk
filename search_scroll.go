package quickwitgosdk

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
)

// SearchScrollRequest is the request body for the Quickwit scroll search API.
type SearchScrollRequest struct {
	SearchRequest
	ScrollTTLSecs uint64 `json:"scroll_ttl_secs,omitempty"`
}

// SearchScroll performs a scroll search and returns a channel of SearchResponse.
// The caller must read from the returned channel until it is closed.
// An error is sent on the errCh channel; only one value is ever sent.
//
// Canceling the context stops reading from the server connection.
func (c *Client) SearchScroll(ctx context.Context, indexId string, req SearchScrollRequest) (<-chan SearchResponse, <-chan error) {
	respCh := make(chan SearchResponse)
	errCh := make(chan error, 1)

	go func() {
		defer close(respCh)
		defer close(errCh)

		body, err := json.Marshal(req)
		if err != nil {
			errCh <- err
			return
		}

		url := c.client.BaseURL + "/api/v1/" + indexId + "/search/scroll"
		httpReq, err := newRequest(ctx, c.client, "POST", url, bytes.NewReader(body))
		if err != nil {
			errCh <- err
			return
		}

		httpResp, err := c.client.GetClient().Do(httpReq)
		if err != nil {
			errCh <- err
			return
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode >= 400 {
			var qe QuickwitError
			if decodeErr := json.NewDecoder(httpResp.Body).Decode(&qe); decodeErr == nil {
				qe.StatusCode = httpResp.StatusCode
				errCh <- &qe
			} else {
				errCh <- &QuickwitError{StatusCode: httpResp.StatusCode, Message: httpResp.Status}
			}
			return
		}

		scanner := bufio.NewScanner(httpResp.Body)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			var resp SearchResponse
			if decErr := json.Unmarshal(line, &resp); decErr != nil {
				continue
			}
			select {
			case respCh <- resp:
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}
		if scanErr := scanner.Err(); scanErr != nil {
			errCh <- scanErr
		}
	}()

	return respCh, errCh
}

// SearchScrollWithCallback performs a scroll search and calls the callback for each response.
// Canceling the context stops reading from the server connection.
func (c *Client) SearchScrollWithCallback(ctx context.Context, indexId string, req SearchScrollRequest, callback func(SearchResponse) error) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := c.client.BaseURL + "/api/v1/" + indexId + "/search/scroll"
	httpReq, err := newRequest(ctx, c.client, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpResp, err := c.client.GetClient().Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode >= 400 {
		var qe QuickwitError
		if decodeErr := json.NewDecoder(httpResp.Body).Decode(&qe); decodeErr == nil {
			qe.StatusCode = httpResp.StatusCode
			return &qe
		}
		return &QuickwitError{StatusCode: httpResp.StatusCode, Message: httpResp.Status}
	}

	scanner := bufio.NewScanner(httpResp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var resp SearchResponse
		if decErr := json.Unmarshal(line, &resp); decErr != nil {
			continue
		}
		if cbErr := callback(resp); cbErr != nil {
			return cbErr
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	return scanner.Err()
}

