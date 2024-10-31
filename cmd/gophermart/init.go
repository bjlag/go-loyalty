package main

import (
	nativeLog "log"

	"github.com/bjlag/go-loyalty/internal/infrastructure/config"
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
