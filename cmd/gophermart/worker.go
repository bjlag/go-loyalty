package main

import (
	"context"
	"time"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/usecase/accrual/update"
)

type accrualWorker struct {
	usecase *update.Usecase
	log     logger.Logger
}

func newAccrualWorker(usecase *update.Usecase, log logger.Logger) *accrualWorker {
	return &accrualWorker{
		usecase: usecase,
		log:     log,
	}
}

func (w *accrualWorker) run(ctx context.Context) {
	resultCh := make(chan *update.Result)

	go func() {
		w.log.Info("Accrual worker started")

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				w.log.Info("Stopped accrual worker")
				return
			case <-ticker.C:
				err := w.usecase.Update(ctx, resultCh)
				if err != nil {
					w.log.WithError(err).Error("Failed to update accrual")
					continue
				}
			}
		}
	}()

	go func() {
		for result := range resultCh {
			if result == nil {
				continue
			}

			log := w.log.
				WithField("order", result.OrderNumber).
				WithField("user", result.UserGUID).
				WithField("old_status", result.OldStatus.String()).
				WithField("old_accrual", result.OldAccrual).
				WithField("new_status", result.NewStatus.String()).
				WithField("new_accrual", result.NewAccrual)

			if result.Err != nil {
				log.WithError(result.Err).Error("Failed to update accrual")
				continue
			}

			log.Info("Accrual updated")
		}
	}()
}
