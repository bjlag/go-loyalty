package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bjlag/go-loyalty/internal/infrastructure/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	c := config.Parse()
	if c == nil {
		log.Fatal("config is nil")
	}

	app := newApp(
		withServerAddr(c.RunAddr().Host(), c.RunAddr().Port()),
	)

	if err := app.run(ctx); err != nil {
		log.Fatal(err)
	}
}
