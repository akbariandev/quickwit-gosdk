package quickwitgosdk

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func (msg *ErrorMessage) Error() string {
	return fmt.Sprintf("quickwitgo error: %s", msg.Message)
}

func NewClient(baseURL string) *Client {
	httpClient := resty.New().
		SetBaseURL(baseURL).
		OnAfterResponse(func(client *resty.Client, resp *resty.Response) error {
			if !resp.IsError() {
				return nil
			}

			if errMsg, ok := resp.Error().(*ErrorMessage); ok {
				return errMsg
			}

			return nil
		})

	return &Client{httpClient}
}
