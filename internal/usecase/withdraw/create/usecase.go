package create

import (
	"context"
	"errors"
	"time"

	"github.com/bjlag/go-loyalty/internal/infrastructure/guid"
	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
	"github.com/bjlag/go-loyalty/internal/model"
)

var (
	ErrInsufficientBalanceOnAccount = errors.New("insufficient balance on account")
)

type Usecase struct {
	accrualRepo repository.AccrualRepo
	accountRepo repository.AccountRepo
	guidGen     guid.IGenerator
}

func NewUsecase(accrualRepo repository.AccrualRepo, accountRepo repository.AccountRepo, guidGen guid.IGenerator) *Usecase {
	return &Usecase{
		accrualRepo: accrualRepo,
		accountRepo: accountRepo,
		guidGen:     guidGen,
	}
}

func (u *Usecase) CreateWithdraw(ctx context.Context, accountGUID, orderNumber string, sum float32) error {
	balance, err := u.accountRepo.Balance(ctx, accountGUID)
	if err != nil {
		return err
	}

	if sum > balance {
		return ErrInsufficientBalanceOnAccount
	}

	transaction := model.NewWithdrawTransaction(
		u.guidGen.Generate(),
		accountGUID,
		orderNumber,
		uint(sum),
		time.Now(),
	)

	err = u.accrualRepo.WithdrawBalance(ctx, transaction)
	if err != nil {
		return err
	}

	return nil
}
