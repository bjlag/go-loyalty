package create

import (
	"context"
	"errors"

	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
	"github.com/bjlag/go-loyalty/internal/model"
)

var (
	ErrAnotherUserHasAlreadyRegisteredOrder = errors.New("another user has already registered an order")
	ErrOrderAlreadyExists                   = errors.New("order already exists")
)

type Usecase struct {
	repo repository.AccrualRepo
}

func NewUsecase(repo repository.AccrualRepo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) CreateAccrual(ctx context.Context, accrual *model.Accrual) error {
	if existAccrual, err := u.repo.AccrualByOrderNumber(ctx, accrual.OrderNumber); err != nil || existAccrual != nil {
		if err != nil {
			return err
		}

		if existAccrual.UserGUID != accrual.UserGUID {
			return ErrAnotherUserHasAlreadyRegisteredOrder
		}

		return ErrOrderAlreadyExists
	}

	err := u.repo.Create(ctx, accrual)
	if err != nil {
		return err
	}

	return nil
}
