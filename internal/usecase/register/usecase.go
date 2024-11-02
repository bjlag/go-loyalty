package register

import (
	"context"
	"errors"
	"fmt"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/guid"
	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
	"github.com/bjlag/go-loyalty/internal/model"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type Usecase struct {
	userRepo *repository.UserRepository
	guidGen  *guid.Generator
	hasher   *auth.Hasher
	jwt      *auth.JWTBuilder
}

func NewUsecase(userRepo *repository.UserRepository, guidGen *guid.Generator, hasher *auth.Hasher, jwt *auth.JWTBuilder) *Usecase {
	return &Usecase{
		userRepo: userRepo,
		guidGen:  guidGen,
		hasher:   hasher,
		jwt:      jwt,
	}
}

func (c Usecase) RegisterUser(ctx context.Context, login, password string) (string, error) {
	user, err := c.userRepo.FindByEmail(ctx, login)
	if err != nil {
		return "", err
	}
	if user != nil {
		return "", fmt.Errorf("%w: email %q", ErrUserAlreadyExists, login)
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return "", err
	}

	user = &model.User{
		GUID:     c.guidGen.Generate(),
		Email:    login,
		Password: hashedPassword,
	}

	err = c.userRepo.Insert(ctx, user)
	if err != nil {
		return "", err
	}

	return c.jwt.BuildJWTString(user.GUID)
}
