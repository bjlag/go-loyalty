package main

import (
	nativeLog "log"
	"os"

	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-loyalty/internal/infrastructure/config"
	"github.com/bjlag/go-loyalty/internal/infrastructure/db/migrator"
	"github.com/bjlag/go-loyalty/internal/infrastructure/db/pg"
	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
)

func mustInitConfig() *config.Configuration {
	cfg := config.Parse()
	if cfg == nil {
		nativeLog.Fatal("config is nil")
	}

	return cfg
}

func mustInitLog(level string) *logger.ZapLog {
	log, err := logger.NewZapLog(level)
	if err != nil {
		nativeLog.Fatalf("faild to create logger: %v", err)
	}

	return log
}

func mustInitDB(dsn string, log logger.Logger) *sqlx.DB {
	db, err := pg.Connect(dsn)
	if err != nil {
		log.WithField("dsn", dsn).
			WithError(err).
			Error("Unable to connect to database")
		os.Exit(1)
	}

	log.WithField("dsn", dsn).Info("Database connection is established")

	return db
}

func mustUpMigrate(source string, db *sqlx.DB, log logger.Logger) {
	driver, err := pgx.WithInstance(db.DB, &pgx.Config{})
	if err != nil {
		log.WithError(err).Error("Error creating database driver")
		os.Exit(1)
	}

	migrate, err := migrator.NewMigrator(source, driver)
	if err != nil {
		log.WithError(err).WithField("source", source).Error("Error initializing migrator")
		os.Exit(1)
	}

	if updated, err := migrate.Up(); err != nil || updated {
		if err != nil {
			log.WithError(err).Error("Migrate error")
			os.Exit(1)
		}
		log.Info("Database migrate is successful")
	}
}
