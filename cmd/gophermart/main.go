package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bjlag/go-loyalty/internal/api/handler/register"
	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	ucRegister "github.com/bjlag/go-loyalty/internal/usecase/register"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := mustInitConfig()
	log := mustInitLog(cfg.LogLevel())
	defer log.Close()

	log.Infof("Log level %q", cfg.LogLevel())
	log.Infof("Run address \"%s:%d\"", cfg.RunAddrHost(), cfg.RunAddrPort())
	log.Infof("JWT secret key %q", cfg.JWTSecretKey())
	log.Infof("JWT expiration time %q", cfg.JWTExpTime())

	hasher := auth.NewHasher()
	jwtBuilder := auth.NewJWTBuilder(cfg.JWTSecretKey(), cfg.JWTExpTime())
	usecaseRegister := ucRegister.NewUsecase(hasher, jwtBuilder)

	app := newApp(
		withRunAddr(cfg.RunAddrHost(), cfg.RunAddrPort()),
		withLogger(log),

		withAPIHandler(http.MethodPost, "/api/user/register", register.NewHandler(usecaseRegister, log).Handle),
	)

	if err := app.run(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("Shut down")
			return
		}

		log.WithError(err).Error("App error")
	}
}
