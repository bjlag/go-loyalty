package model

import "time"

type Account struct {
	GUID        string
	Balance     float64
	WithdrawSum float64
	UpdatedAt   time.Time
}
