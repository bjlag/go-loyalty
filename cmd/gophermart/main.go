package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bjlag/go-loyalty/internal/api/handler/balance/get"
	"github.com/bjlag/go-loyalty/internal/api/handler/balance/withdraw"
	"github.com/bjlag/go-loyalty/internal/api/handler/order/list"
	"github.com/bjlag/go-loyalty/internal/api/handler/order/upload"
	"github.com/bjlag/go-loyalty/internal/api/handler/user/login"
	"github.com/bjlag/go-loyalty/internal/api/handler/user/register"
	"github.com/bjlag/go-loyalty/internal/api/handler/withdrawals"
	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/client"
	"github.com/bjlag/go-loyalty/internal/infrastructure/guid"
	"github.com/bjlag/go-loyalty/internal/infrastructure/middleware"
	"github.com/bjlag/go-loyalty/internal/infrastructure/repository"
	"github.com/bjlag/go-loyalty/internal/infrastructure/service/accrual"
	ucCreateAccrual "github.com/bjlag/go-loyalty/internal/usecase/accrual/create"
	ucUpdateAccrual "github.com/bjlag/go-loyalty/internal/usecase/accrual/update"
	ucLogin "github.com/bjlag/go-loyalty/internal/usecase/user/login"
	ucRegister "github.com/bjlag/go-loyalty/internal/usecase/user/register"
	ucCreateWithdraw "github.com/bjlag/go-loyalty/internal/usecase/withdraw/create"
)

const (
	accrualTimeout       = 200 * time.Millisecond
	accrualRetryCount    = 2
	accrualRetryWaitTime = 100 * time.Millisecond
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := mustInitConfig()
	log := mustInitLog(cfg.LogLevel())
	defer log.Close()

	log.Infof("Log level %q", cfg.LogLevel())
	log.Infof("Run address \"%s:%d\"", cfg.RunAddrHost(), cfg.RunAddrPort())
	log.Infof("Accrual address %q", cfg.AccrualSystemAddress())
	log.Infof("JWT secret key %q", cfg.JWTSecretKey())
	log.Infof("JWT expiration time %q", cfg.JWTExpTime())
	log.Infof("Database URI %q", cfg.DatabaseURI())
	log.Infof("Path to migration source files %q", cfg.MigratePath())

	db := mustInitDB(cfg.DatabaseURI(), log)
	mustUpMigrate(cfg.MigratePath(), db, log)

	userRepo := repository.NewUserPG(db)
	accrualRepo := repository.NewAccrualPG(db)
	accountRepo := repository.NewAccountPG(db)
	transactionRepo := repository.NewTransactionPG(db)

	hasher := auth.NewHasher()
	jwtBuilder := auth.NewJWTBuilder(cfg.JWTSecretKey(), cfg.JWTExpTime())

	accrualClient := accrual.NewAccrualClient(
		client.NewRestyClient(
			client.WithTimeout(accrualTimeout),
			client.WithRetryCount(accrualRetryCount),
			client.WithRetryWaitTime(accrualRetryWaitTime),
			client.WithLogger(log),
		),
		cfg.AccrualSystemAddress(),
	)

	guidGen := new(guid.Generator)

	usecaseRegister := ucRegister.NewUsecase(userRepo, guidGen, hasher, jwtBuilder)
	usecaseLogin := ucLogin.NewUsecase(userRepo, hasher, jwtBuilder)
	usecaseCreateAccrual := ucCreateAccrual.NewUsecase(accrualRepo)
	usecaseUpdateAccrual := ucUpdateAccrual.NewUsecase(accrualClient, accrualRepo, guidGen)
	usecaseCreateWithdraw := ucCreateWithdraw.NewUsecase(accrualRepo, accountRepo, guidGen)

	worker := newAccrualWorker(usecaseUpdateAccrual, log)
	worker.run(ctx)

	app := newApp(
		withRunAddr(cfg.RunAddrHost(), cfg.RunAddrPort()),
		withLogger(log),

		withAPIHandler(http.MethodPost, "/api/user/register", register.NewHandler(usecaseRegister, log).Handle),
		withAPIHandler(http.MethodPost, "/api/user/login", login.NewHandler(usecaseLogin, log).Handle),

		withAPIHandler(http.MethodPost, "/api/user/orders", upload.NewHandler(usecaseCreateAccrual, log).Handle, middleware.CheckAuth(jwtBuilder, log)),
		withAPIHandler(http.MethodGet, "/api/user/orders", list.NewHandler(accrualRepo, log).Handle, middleware.CheckAuth(jwtBuilder, log)),

		withAPIHandler(http.MethodGet, "/api/user/balance", get.NewHandler(accountRepo, log).Handle, middleware.CheckAuth(jwtBuilder, log)),
		withAPIHandler(http.MethodPost, "/api/user/balance/withdraw", withdraw.NewHandler(usecaseCreateWithdraw, log).Handle, middleware.CheckAuth(jwtBuilder, log)),
		withAPIHandler(http.MethodGet, "/api/user/withdrawals", withdrawals.NewHandler(transactionRepo, log).Handle, middleware.CheckAuth(jwtBuilder, log)),
	)

	if err := app.run(ctx); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("Shut down")
			return
		}

		log.WithError(err).Error("App error")
	}
}
