package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
	serviceAccrual "github.com/bjlag/go-loyalty/internal/infrastructure/service/accrual"
	"github.com/bjlag/go-loyalty/internal/model"
)

var (
	mapAccrualStatus = map[string]model.AccrualStatus{
		"registered": model.New,
		"processing": model.Processing,
		"invalid":    model.Invalid,
		"processed":  model.Processed,
	}
)

type Usecase struct {
	client *serviceAccrual.Client
	repo   repository.AccrualRepository
}

func NewUsecase(client *serviceAccrual.Client, repo repository.AccrualRepository) *Usecase {
	return &Usecase{
		client: client,
		repo:   repo,
	}
}

func (u Usecase) Update(ctx context.Context) error {
	accrualsInWork, err := u.repo.AccrualsInWork(ctx)
	if err != nil {
		return err
	}

	for _, accrual := range accrualsInWork {
		resp, err := u.client.OrderStatus(accrual.OrderNumber)
		if err != nil {
			continue
		}

		newStatus, ok := mapAccrualStatus[strings.ToLower(resp.Status)]
		if !ok {
			continue
		}

		if newStatus == accrual.Status {
			continue
		}

		var newAccrual uint
		if resp.Accrual != nil {
			newAccrual = *resp.Accrual
		}

		err = u.repo.Update(ctx, accrual.OrderNumber, newStatus, newAccrual)
		if err != nil {
			return err
		}

		fmt.Println(resp)
	}

	return nil
}
