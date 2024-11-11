package model

import "time"

type AccrualStatus uint

const (
	New AccrualStatus = iota
	Processing
	Invalid
	Processed
)

func (s AccrualStatus) String() string {
	switch s {
	case New:
		return "New"
	case Processing:
		return "Processing"
	case Invalid:
		return "Invalid"
	case Processed:
		return "Processed"
	}
	return "Unknown"
}

type Accrual struct {
	OrderNumber string
	UserGUID    string
	Status      AccrualStatus
	Accrual     uint
	UploadedAt  time.Time
}

func NewAccrual(orderNumber, userGUID string) *Accrual {
	return &Accrual{
		OrderNumber: orderNumber,
		UserGUID:    userGUID,
		Status:      New,
		UploadedAt:  time.Now(),
	}
}
