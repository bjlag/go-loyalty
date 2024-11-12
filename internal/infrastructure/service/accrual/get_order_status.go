package accrual

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual *int64 `json:"accrual,omitempty"`
}

func (c Client) OrderStatus(orderNumber string) (*Response, error) {
	resp, err := c.client.Get(c.serviceURL + "/api/orders/" + orderNumber)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	var result *Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
	}

	return result, nil
}
