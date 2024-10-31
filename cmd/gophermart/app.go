package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

const (
	serverHost = "localhost"
	serverPort = 8080
)

type option func(a *application)

func withServerHost(host string) option {
	return func(a *application) {
		a.serverHost = host
	}
}

func withServerPort(port int) option {
	return func(a *application) {
		a.serverPort = port
	}
}

type application struct {
	serverHost string
	serverPort int
}

func newApp(opts ...option) *application {
	a := &application{}
	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a application) run(ctx context.Context) error {
	if a.serverHost == "" {
		a.serverHost = serverHost
	}

	if a.serverPort == 0 {
		a.serverPort = serverPort
	}

	log.Printf("Starting server on %s:%d", a.serverHost, a.serverPort)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", a.serverHost, a.serverPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, world!"))
		}),
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return server.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()

		log.Println("Graceful shutting down server")

		return server.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
