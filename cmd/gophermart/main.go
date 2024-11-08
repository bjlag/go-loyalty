package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bjlag/go-loyalty/internal/api/handler/login"
	"github.com/bjlag/go-loyalty/internal/api/handler/register"
	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/guid"
	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
	ucLogin "github.com/bjlag/go-loyalty/internal/usecase/login"
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
	log.Infof("Database URI %q", cfg.DatabaseURI())
	log.Infof("Path to migration source files %q", cfg.MigratePath())

	db := mustInitDB(cfg.DatabaseURI(), log)
	mustUpMigrate(cfg.MigratePath(), db, log)

	userRepo := repository.NewUserRepository(db)
	hasher := auth.NewHasher()
	jwtBuilder := auth.NewJWTBuilder(cfg.JWTSecretKey(), cfg.JWTExpTime())
	usecaseRegister := ucRegister.NewUsecase(userRepo, new(guid.Generator), hasher, jwtBuilder)
	usecaseLogin := ucLogin.NewUsecase(userRepo, hasher, jwtBuilder)

	app := newApp(
		withRunAddr(cfg.RunAddrHost(), cfg.RunAddrPort()),
		withLogger(log),

		withAPIHandler(http.MethodPost, "/api/user/register", register.NewHandler(usecaseRegister, log).Handle),
		withAPIHandler(http.MethodPost, "/api/user/login", login.NewHandler(usecaseLogin, log).Handle),
	)

	if err := app.run(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("Shut down")
			return
		}

		log.WithError(err).Error("App error")
	}
}
