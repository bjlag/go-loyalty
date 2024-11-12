package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/infrastructure/service/accrual"
)

type accrualWorker struct {
	client *accrual.Client
	log    logger.Logger
}

func newAccrualWorker(client *accrual.Client, log logger.Logger) *accrualWorker {
	return &accrualWorker{
		client: client,
		log:    log,
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
				resp, err := w.client.OrderStatus("12345678705")
				if err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println(resp)
			}
		}
	}()
}
