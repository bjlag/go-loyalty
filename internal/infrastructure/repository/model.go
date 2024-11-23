package repository

import (
	"time"

	"github.com/bjlag/go-loyalty/internal/model"
)

type user struct {
	GUID     string `db:"guid"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

func accrualFromModel(model model.User) *user {
	return &user{
		GUID:     model.GUID,
		Login:    model.Login,
		Password: model.Password,
	}
}

func (u user) export() *model.User {
	return &model.User{
		GUID:     u.GUID,
		Login:    u.Login,
		Password: u.Password,
	}
}

type accrual struct {
	OrderNumber string    `db:"order_number"`
	UserGUID    string    `db:"user_guid"`
	Status      float64   `db:"status"`
	Accrual     float64   `db:"accrual"`
	UploadedAt  time.Time `db:"uploaded_at"`
}

func (a accrual) export() *model.Accrual {
	return &model.Accrual{
		OrderNumber: a.OrderNumber,
		UserGUID:    a.UserGUID,
		Status:      model.AccrualStatus(a.Status),
		Accrual:     a.Accrual,
		UploadedAt:  a.UploadedAt,
	}
}

type transaction struct {
	GUID        string                `db:"guid"`
	AccountGUID string                `db:"account_guid"`
	OrderNumber string                `db:"order_number"`
	Type        model.TransactionType `db:"type"`
	Sum         float64               `db:"sum"`
	ProcessedAt time.Time             `db:"processed_at"`
}

func (t transaction) export() *model.Transaction {
	return &model.Transaction{
		GUID:        t.GUID,
		AccountGUID: t.AccountGUID,
		OrderNumber: t.OrderNumber,
		Type:        t.Type,
		Sum:         t.Sum,
		ProcessedAt: t.ProcessedAt,
	}
}
