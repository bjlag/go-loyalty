package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app := newApp(
		withServerHost("localhost"),
		withServerPort(8080),
	)

	if err := app.run(ctx); err != nil {
		log.Fatal(err)
	}
}
