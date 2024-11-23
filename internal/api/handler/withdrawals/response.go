package withdrawals

import "github.com/bjlag/go-loyalty/internal/api"

type Response []Withdraw

type Withdraw struct {
	Order       string       `json:"order"`
	Sum         float64      `json:"sum"`
	ProcessedAt api.Datetime `json:"processed_at"`
}
