package update

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/sync/errgroup"

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
	repo   repository.AccrualRepo
}

type Result struct {
	Err error
}

func NewResult(err error) *Result {
	return &Result{
		Err: err,
	}
}

func NewUsecase(client *serviceAccrual.Client, repo repository.AccrualRepo) *Usecase {
	return &Usecase{
		client: client,
		repo:   repo,
	}
}

func (u Usecase) Update(ctx context.Context, resultCh chan *Result) error {
	var err error

	accrualsInWork, err := u.repo.AccrualsInWork(ctx)
	if err != nil {
		return err
	}

	g, gCtx := errgroup.WithContext(ctx)

	for _, accrual := range accrualsInWork {
		g.Go(func() error {
			select {
			case <-gCtx.Done():
				close(resultCh)
				return gCtx.Err()
			default:
			}

			resp, err := u.client.OrderStatus(accrual.OrderNumber)
			if err != nil {
				resultCh <- NewResult(err)
				return nil
			}

			newStatus, ok := mapAccrualStatus[strings.ToLower(resp.Status)]
			if !ok {
				resultCh <- NewResult(fmt.Errorf("unknown status: %s", resp.Status))
				return nil
			}

			if newStatus == accrual.Status {
				return nil
			}

			var newAccrual uint
			if resp.Accrual != nil {
				newAccrual = *resp.Accrual
			}

			if newAccrual > 0 {
				err = u.repo.AddTx(gCtx, model.Accrual{
					OrderNumber: accrual.OrderNumber,
					UserGUID:    accrual.UserGUID,
					Status:      newStatus,
					Accrual:     newAccrual,
					UploadedAt:  accrual.UploadedAt,
				})
			} else {
				err = u.repo.UpdateStatus(gCtx, accrual.OrderNumber, newStatus)
			}

			if err != nil {
				resultCh <- NewResult(err)
				return nil
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
