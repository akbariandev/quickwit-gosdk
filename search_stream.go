package quickwitgosdk

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"

	"net/http"

	"github.com/go-resty/resty/v2"
)

// SearchStreamRequest is the request body for the Quickwit search stream API.
type SearchStreamRequest struct {
	SearchRequest
	FastField    string `json:"fast_field,omitempty"`
	OutputFormat string `json:"output_format,omitempty"` // "csv" or "click_house_row_binary"
}

// SearchStream performs a streaming search and returns a channel of SearchResponse.
// The caller must read from the returned channel until it is closed.
// An error is sent on the errCh channel; only one value is ever sent.
//
// Canceling the context stops reading from the server connection.
func (c *Client) SearchStream(ctx context.Context, indexId string, req SearchStreamRequest) (<-chan SearchResponse, <-chan error) {
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

		// Use the underlying http.Client directly for streaming — resty buffers the full response.
		url := c.client.BaseURL + "/api/v1/" + indexId + "/search/stream"
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

		// Read NDJSON lines from the response body.
		scanner := bufio.NewScanner(httpResp.Body)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			var resp SearchResponse
			if decErr := json.Unmarshal(line, &resp); decErr != nil {
				// Best-effort: skip malformed lines.
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

// SearchStreamWithCallback performs a streaming search and calls the callback for each response.
// Canceling the context stops reading from the server connection.
func (c *Client) SearchStreamWithCallback(ctx context.Context, indexId string, req SearchStreamRequest, callback func(SearchResponse) error) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := c.client.BaseURL + "/api/v1/" + indexId + "/search/stream"
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

// newRequest builds an *http.Request from the resty client's settings.
func newRequest(ctx context.Context, rc *resty.Client, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Carry over resty default headers.
	if rc.Token != "" {
		req.Header.Set("Authorization", rc.AuthScheme+" "+rc.Token)
	}
	for k, vv := range rc.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
