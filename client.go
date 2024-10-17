package quickwitgosdk

import "github.com/go-resty/resty/v2"

type Client struct {
	client *resty.Client
}

func NewClient(baseURL string) *Client {
	httpClient := resty.New().SetBaseURL(baseURL)
	return &Client{httpClient}
}
