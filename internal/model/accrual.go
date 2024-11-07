package model

import "time"

type AccrualStatus uint

const (
	New AccrualStatus = iota
	Processing
	Invalid
	Processed
)

type Accrual struct {
	Number     string
	UserGUID   string
	Status     AccrualStatus
	Accrual    uint
	UploadedAt time.Time
}

func NewAccrual(number, userGUID string) *Accrual {
	return &Accrual{
		Number:     number,
		UserGUID:   userGUID,
		Status:     New,
		UploadedAt: time.Now(),
	}
}
