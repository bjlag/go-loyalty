package repository

import "github.com/bjlag/go-loyalty/internal/model"

type user struct {
	GUID     string `db:"guid"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func userFromModel(model model.User) *user {
	return &user{
		GUID:     model.GUID,
		Email:    model.Email,
		Password: model.Password,
	}
}

func (u user) export() *model.User {
	return &model.User{
		GUID:     u.GUID,
		Email:    u.Email,
		Password: u.Password,
	}
}
