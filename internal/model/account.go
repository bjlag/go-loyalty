package model

import "time"

type Account struct {
	GUID        string
	Balance     uint
	WithdrawSum uint
	UpdatedAt   time.Time
}
