package quickwitgosdk

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// ClientOption is a functional option for configuring the Quickwit client.
type ClientOption func(*Client)

// WithAPIKey sets the API key used for Bearer token authentication.
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.client.SetAuthToken(apiKey)
	}
}

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
	c := &Client{client: httpClient}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
