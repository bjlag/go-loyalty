package model

import "time"

type Account struct {
	GUID      string
	Balance   uint
	UpdatedAt time.Time
}
