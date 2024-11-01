package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := mustInitConfig()
	log := mustInitLog(cfg.LogLevel())
	defer log.Close()

	log.Infof("Log level %q", cfg.LogLevel())

	addr := runAddr{host: cfg.RunAddrHost(), port: cfg.RunAddrPort()}
	app := newApp(addr, log)

	if err := app.run(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("Shut down")
			return
		}

		log.WithError(err).Error("App error")
	}
}
