package list

import (
	"github.com/bjlag/go-loyalty/internal/api"
)

type Response []Order

type Order struct {
	Number     string       `json:"number"`
	Status     string       `json:"status"`
	Accrual    uint         `json:"accrual"`
	UploadedAt api.Datetime `json:"uploaded_at"`
}
