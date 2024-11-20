package model

import "time"

type Transaction struct {
	GUID        string
	AccountGUID string
	OrderNumber string
	Sum         int
	ProcessedAt time.Time
}
