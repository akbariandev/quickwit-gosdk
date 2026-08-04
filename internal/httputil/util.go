// Package httputil provides shared HTTP utilities for the Quickwit Go SDK.
package httputil

import (
	"encoding/json"
	"strings"

	"github.com/akbariandev/quickwit-gosdk/types"
	"github.com/go-resty/resty/v2"
)

// SetupErrorHook registers an OnAfterResponse hook on the resty client that
// converts any non-2xx response into a types.Error, extracting Quickwit's
// error message from the JSON body (e.g. {"message": "..."}).
func SetupErrorHook(c *resty.Client) {
	c.OnAfterResponse(func(_ *resty.Client, resp *resty.Response) error {
		if resp.IsSuccess() {
			return nil
		}

		var body struct {
			Message string `json:"message"`
		}
		msg := resp.Status()
		if err := json.Unmarshal(resp.Body(), &body); err == nil && strings.TrimSpace(body.Message) != "" {
			msg = body.Message
		}
		return &types.Error{StatusCode: resp.StatusCode(), Message: msg}
	})
}

// BoolToString converts a bool to its string representation for query parameters.
func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
