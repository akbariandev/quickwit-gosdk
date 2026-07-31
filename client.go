package quickwitgosdk

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// ClientOption is a functional option for configuring the Quickwit client.
type ClientOption func(*Client)

// WithTimeout sets the HTTP request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.client.SetTimeout(timeout)
	}
}

// WithTransport sets a custom http.RoundTripper on the underlying HTTP client (e.g. for custom TLS or logging).
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(c *Client) {
		c.client.SetTransport(transport)
	}
}

// Client is a Quickwit API client.
type Client struct {
	client *resty.Client
}

// NewClient creates a new Quickwit client with the given base URL and optional configuration.
func NewClient(baseURL string, opts ...ClientOption) *Client {
	httpClient := resty.New().SetBaseURL(baseURL)

	// Treat any non-2xx response as an error. Quickwit returns error details as JSON
	// in the body (e.g. {"message": "..."}); extract and surface them as QuickwitError.
	httpClient.OnAfterResponse(func(_ *resty.Client, resp *resty.Response) error {
		if resp.IsSuccess() {
			return nil
		}

		// Try to extract Quickwit's error message from the JSON body.
		var body struct {
			Message string `json:"message"`
		}
		msg := resp.Status()
		if err := json.Unmarshal(resp.Body(), &body); err == nil && strings.TrimSpace(body.Message) != "" {
			msg = body.Message
		}
		return &QuickwitError{StatusCode: resp.StatusCode(), Message: msg}
	})

	c := &Client{client: httpClient}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
