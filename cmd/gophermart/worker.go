package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bjlag/go-loyalty/internal/infrastructure/client"
	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/infrastructure/service/accrual"
)

type accrualWorker struct {
	accrualAddr addr
	log         logger.Logger
}

func newAccrualWorker(accrualAddrHost string, accrualAddrPort int, log logger.Logger) *accrualWorker {
	return &accrualWorker{
		accrualAddr: addr{
			host: accrualAddrHost,
			port: accrualAddrPort,
		},
		log: log,
	}
}

func (w *accrualWorker) run(ctx context.Context) {
	go func() {
		w.log.Info("Accrual worker started")

		restyClient := client.NewRestyClient(
			client.WithTimeout(200*time.Millisecond),
			client.WithRetryCount(2),
			client.WithRetryWaitTime(100*time.Millisecond),
		)

		accrualClient := accrual.NewAccrualClient(restyClient, w.accrualAddr.host, w.accrualAddr.port)

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				w.log.Info("Stopped accrual worker")
				return
			case <-ticker.C:
				resp, err := accrualClient.OrderStatus("12345678705")
				if err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println(resp)
			}
		}
	}()
}
