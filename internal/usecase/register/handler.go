package register

import (
	"context"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/model"
)

type Usecase struct {
	hasher *auth.Hasher
	jwt    *auth.JWTBuilder
}

func NewUsecase(hasher *auth.Hasher, jwt *auth.JWTBuilder) *Usecase {
	return &Usecase{
		hasher: hasher,
		jwt:    jwt,
	}
}

func (c Usecase) RegisterUser(ctx context.Context, login, password string) (string, error) {
	// ищем пользователя с переданным логином, если есть, то ошибка

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &model.User{
		GUID:     "guid",
		Email:    login,
		Password: hashedPassword,
	}
	// регистрируем пользователя

	return c.jwt.BuildJWTString(user.GUID)
}
