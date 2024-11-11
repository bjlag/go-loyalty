package list

import "time"

type Response []Order

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    uint      `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}
