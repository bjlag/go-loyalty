package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)

type Usecase struct {
	userRepo *repository.UserRepository
	hasher   *auth.Hasher
	jwt      *auth.JWTBuilder
}

func NewUsecase(userRepo *repository.UserRepository, hasher *auth.Hasher, jwt *auth.JWTBuilder) *Usecase {
	return &Usecase{
		userRepo: userRepo,
		hasher:   hasher,
		jwt:      jwt,
	}
}

func (u *Usecase) LoginUser(ctx context.Context, login, password string) (string, error) {
	user, err := u.userRepo.FindByEmail(ctx, login)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("%w: email %q", ErrUserNotFound, login)
	}

	if !u.hasher.ComparePasswords(user.Password, password) {
		return "", fmt.Errorf("%w: email %q", ErrWrongPassword, login)
	}

	return u.jwt.BuildJWTString(user.GUID)
}
