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
	Status      uint      `db:"status"`
	Accrual     uint      `db:"accrual"`
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
