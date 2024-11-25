package model

import "time"

type TransactionType uint

const (
	Add      TransactionType = iota // Начислили на счет
	Withdraw                        // Сняли со счета
)

type Transaction struct {
	GUID        string
	AccountGUID string
	OrderNumber string
	Type        TransactionType
	Sum         float64
	ProcessedAt time.Time
}

func NewAddTransaction(guid, accountGUID, orderNumber string, sum float64, processedAt time.Time) Transaction {
	return Transaction{
		GUID:        guid,
		AccountGUID: accountGUID,
		OrderNumber: orderNumber,
		Type:        Add,
		Sum:         sum,
		ProcessedAt: processedAt,
	}
}

func NewWithdrawTransaction(guid, accountGUID, orderNumber string, sum float64, processedAt time.Time) Transaction {
	return Transaction{
		GUID:        guid,
		AccountGUID: accountGUID,
		OrderNumber: orderNumber,
		Type:        Withdraw,
		Sum:         sum,
		ProcessedAt: processedAt,
	}
}
