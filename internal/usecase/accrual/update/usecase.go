package update

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/bjlag/go-loyalty/internal/infrastructure/guid"
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
	client  *serviceAccrual.Client
	repo    repository.AccrualRepo
	guidGen guid.IGenerator
}

type Result struct {
	OrderNumber string
	UserGUID    string
	OldStatus   model.AccrualStatus
	OldAccrual  float64
	NewStatus   *model.AccrualStatus
	NewAccrual  *float64
	Err         error
}

func NewResult(
	orderNumber string,
	userGUID string,
	oldStatus model.AccrualStatus,
	oldAccrual float64,
	newStatus *model.AccrualStatus,
	newAccrual *float64,
	err error,
) *Result {
	return &Result{
		OrderNumber: orderNumber,
		UserGUID:    userGUID,
		OldStatus:   oldStatus,
		OldAccrual:  oldAccrual,
		NewStatus:   newStatus,
		NewAccrual:  newAccrual,
		Err:         err,
	}
}

func NewUsecase(client *serviceAccrual.Client, repo repository.AccrualRepo, guidGen guid.IGenerator) *Usecase {
	return &Usecase{
		client:  client,
		repo:    repo,
		guidGen: guidGen,
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
				resultCh <- NewResult(accrual.OrderNumber, accrual.UserGUID, accrual.Status, accrual.Accrual, nil, nil, err)
				return nil
			}

			newStatus, ok := mapAccrualStatus[strings.ToLower(resp.Status)]
			if !ok {
				resultCh <- NewResult(
					accrual.OrderNumber,
					accrual.UserGUID,
					accrual.Status,
					accrual.Accrual,
					nil,
					nil,
					fmt.Errorf("unknown status: %s", resp.Status),
				)
				return nil
			}

			if newStatus == accrual.Status {
				return nil
			}

			var newAccrual float64
			if resp.Accrual != nil {
				newAccrual = *resp.Accrual
			}

			if newAccrual > 0 {
				mAccrual := model.Accrual{
					OrderNumber: accrual.OrderNumber,
					UserGUID:    accrual.UserGUID,
					Status:      newStatus,
					Accrual:     newAccrual,
					UploadedAt:  accrual.UploadedAt,
				}

				mAccount := model.Account{
					GUID:      accrual.UserGUID,
					Balance:   newAccrual,
					UpdatedAt: time.Now(),
				}

				mTransaction := model.NewAddTransaction(
					u.guidGen.Generate(),
					mAccount.GUID,
					accrual.OrderNumber,
					newAccrual,
					time.Now(),
				)

				err = u.repo.AddBalance(gCtx, mAccrual, mAccount, mTransaction)
			} else {
				err = u.repo.UpdateStatus(gCtx, accrual.OrderNumber, newStatus)
			}

			if err != nil {
				resultCh <- NewResult(accrual.OrderNumber, accrual.UserGUID, accrual.Status, accrual.Accrual, nil, nil, err)
				return nil
			}

			resultCh <- NewResult(accrual.OrderNumber, accrual.UserGUID, accrual.Status, accrual.Accrual, &newStatus, &newAccrual, nil)

			return nil
		})
	}

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
