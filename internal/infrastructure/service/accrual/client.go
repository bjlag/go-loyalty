package accrual

import (
	"fmt"

	"github.com/bjlag/go-loyalty/internal/infrastructure/client"
)

type Client struct {
	client     client.Client
	serviceURL string
}

func NewAccrualClient(client client.Client, host string, port int) *Client {
	return &Client{
		client:     client,
		serviceURL: fmt.Sprintf("http://%s:%d", host, port),
	}
}
