package accrual

import (
	"github.com/bjlag/go-loyalty/internal/infrastructure/client"
)

type Client struct {
	client     client.Client
	serviceURL string
}

func NewAccrualClient(client client.Client, serviceUrl string) *Client {
	return &Client{
		client:     client,
		serviceURL: serviceUrl,
	}
}
