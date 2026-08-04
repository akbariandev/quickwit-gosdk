// Package client provides the Quickwit HTTP client used by the SDK operation packages.
package client

import (
	"net/http"
	"time"

	"github.com/akbariandev/quickwit-gosdk/internal/httputil"
	"github.com/go-resty/resty/v2"
)

// Option is a functional option for configuring the Quickwit client.
type Option func(*Client)

// WithTimeout sets the HTTP request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.HTTP.SetTimeout(timeout)
	}
}

// WithTransport sets a custom http.RoundTripper on the underlying HTTP client
// (e.g. for custom TLS or logging).
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		c.HTTP.SetTransport(transport)
	}
}

// Client is a Quickwit API client. The underlying resty client is exported as
// HTTP so that the operation packages (search, index, ingest, delete) can build
// requests against it.
type Client struct {
	HTTP *resty.Client
}

// New creates a new Quickwit client with the given base URL and optional configuration.
func New(baseURL string, opts ...Option) *Client {
	restyClient := resty.New().SetBaseURL(baseURL)
	httputil.SetupErrorHook(restyClient)

	c := &Client{HTTP: restyClient}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
