package repository

import "github.com/bjlag/go-loyalty/internal/model"

type user struct {
	GUID     string `db:"guid"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

func userFromModel(model model.User) *user {
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
