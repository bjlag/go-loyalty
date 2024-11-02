package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bjlag/go-loyalty/internal/api/handler/register"
	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	register2 "github.com/bjlag/go-loyalty/internal/usecase/register"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := mustInitConfig()
	log := mustInitLog(cfg.LogLevel())
	defer log.Close()

	log.Infof("Log level %q", cfg.LogLevel())

	hasher := auth.NewHasher()
	jwtBuilder := auth.NewJWTBuilder("secret", time.Hour*3)
	usecase := register2.NewUsecase(hasher, jwtBuilder)

	app := newApp(
		withRunAddr(cfg.RunAddrHost(), cfg.RunAddrPort()),
		withLogger(log),

		withAPIHandler(http.MethodPost, "/api/user/register", register.NewHandler(usecase, log).Handle),
	)

	if err := app.run(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("Shut down")
			return
		}

		log.WithError(err).Error("App error")
	}
}
