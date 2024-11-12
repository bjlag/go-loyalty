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
				err := w.usecase.Update(ctx)
				if err != nil {
					w.log.WithError(err).Error("Failed to update accrual")
					continue
				}
			}
		}
	}()
}
