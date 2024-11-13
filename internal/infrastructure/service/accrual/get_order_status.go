package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrOrderNotRegistered = errors.New("order not registered")
	ErrTooManyRequests    = errors.New("too many requests")
	ErrUnknownStatus      = errors.New("unknown status")
)

type Response struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual *uint  `json:"accrual,omitempty"`
}

func (c Client) OrderStatus(orderNumber string) (*Response, error) {
	resp, err := c.client.Get(c.serviceURL + "/api/orders/" + orderNumber)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNoContent:
			return nil, fmt.Errorf("%w: %s", ErrOrderNotRegistered, orderNumber)
		case http.StatusTooManyRequests:
			return nil, fmt.Errorf("%w: %s", ErrTooManyRequests, orderNumber)
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnknownStatus, orderNumber)
		}
	}

	var result *Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
