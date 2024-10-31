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

func withServerAddr(host string, port int) option {
	return func(a *application) {
		a.serverAddr = &serverAddr{
			host: host,
			port: port,
		}
	}
}

type serverAddr struct {
	host string
	port int
}

type application struct {
	serverAddr *serverAddr
}

func newApp(opts ...option) *application {
	a := &application{}
	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a application) run(ctx context.Context) error {
	if a.serverAddr == nil {
		a.serverAddr = &serverAddr{
			host: serverHost,
			port: serverPort,
		}
	} else {
		if a.serverAddr.host == "" {
			a.serverAddr.host = serverHost
		}

		if a.serverAddr.port == 0 {
			a.serverAddr.port = serverPort
		}
	}

	log.Printf("Starting server on %s:%d", a.serverAddr.host, a.serverAddr.port)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", a.serverAddr.host, a.serverAddr.port),
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
