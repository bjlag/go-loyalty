package client

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client interface {
	Get(url string) (*http.Response, error)
}

type Option func(*resty.Client)

func WithTimeout(timeout time.Duration) Option {
	return func(client *resty.Client) {
		client.SetTimeout(timeout)
	}
}

func WithRetryCount(count int) Option {
	return func(client *resty.Client) {
		client.SetRetryCount(count)
	}
}

func WithRetryWaitTime(waitTime time.Duration) Option {
	return func(client *resty.Client) {
		client.SetRetryWaitTime(waitTime)
	}
}

type RestyClient struct {
	client *resty.Client
}

func NewRestyClient(opts ...Option) *RestyClient {
	r := resty.New()
	for _, opt := range opts {
		opt(r)
	}

	return &RestyClient{
		client: r,
	}
}

func (c RestyClient) Get(url string) (*http.Response, error) {
	resp, err := c.client.R().SetDoNotParseResponse(true).Get(url)
	if err != nil {
		return nil, err
	}

	return resp.RawResponse, nil
}
