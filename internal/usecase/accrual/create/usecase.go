package create

import (
	"context"
	"errors"

	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
	"github.com/bjlag/go-loyalty/internal/infrastructure/validator"
	"github.com/bjlag/go-loyalty/internal/model"
)

var (
	ErrInvalidOrderNumber                   = errors.New("invalid order number")
	ErrAnotherUserHasAlreadyRegisteredOrder = errors.New("another user has already registered an order")
	ErrOrderAlreadyExists                   = errors.New("order already exists")
)

type Usecase struct {
	repo repository.AccrualRepository
}

func NewUsecase(repo repository.AccrualRepository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) CreateAccrual(ctx context.Context, userGUID, orderNumber string) error {
	if !validator.CheckLuhn(orderNumber) {
		return ErrInvalidOrderNumber
	}

	if accrual, err := u.repo.AccrualByOrderNumber(ctx, orderNumber); err != nil || accrual != nil {
		if err != nil {
			return err
		}

		if accrual.UserGUID != userGUID {
			return ErrAnotherUserHasAlreadyRegisteredOrder
		}

		return ErrOrderAlreadyExists
	}

	accrual := model.NewAccrual(orderNumber, userGUID)
	err := u.repo.Insert(ctx, accrual)
	if err != nil {
		return err
	}

	return nil
}
