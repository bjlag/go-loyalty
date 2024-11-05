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
	userRepo repository.UserRepository
	guidGen  guid.IGenerator
	hasher   auth.IHasher
	jwt      *auth.JWTBuilder
}

func NewUsecase(userRepo repository.UserRepository, guidGen guid.IGenerator, hasher auth.IHasher, jwt *auth.JWTBuilder) *Usecase {
	return &Usecase{
		userRepo: userRepo,
		guidGen:  guidGen,
		hasher:   hasher,
		jwt:      jwt,
	}
}

func (u Usecase) RegisterUser(ctx context.Context, login, password string) (string, error) {
	user, err := u.userRepo.FindByEmail(ctx, login)
	if err != nil {
		return "", err
	}
	if user != nil {
		return "", fmt.Errorf("%w: email %q", ErrUserAlreadyExists, login)
	}

	hashedPassword, err := u.hasher.HashPassword(password)
	if err != nil {
		return "", err
	}

	user = &model.User{
		GUID:     u.guidGen.Generate(),
		Email:    login,
		Password: hashedPassword,
	}

	err = u.userRepo.Insert(ctx, user)
	if err != nil {
		return "", err
	}

	return u.jwt.BuildJWTString(user.GUID)
}
