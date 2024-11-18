package model

import "time"

type Transaction struct {
	GUID        string
	AccountGUID string
	OrderNumber string
	Sum         uint
	ProcessedAt time.Time
}
