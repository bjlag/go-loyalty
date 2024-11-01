package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
)

type option func(a *application)

type runAddr struct {
	host string
	port int
}

type application struct {
	runAddr runAddr
	log     logger.Logger
}

func newApp(runAddr runAddr, log logger.Logger, opts ...option) *application {
	a := &application{
		runAddr: runAddr,
		log:     log,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a application) run(ctx context.Context) error {
	a.log.
		WithField("host", a.runAddr.host).
		WithField("port", a.runAddr.port).
		Info("Starting server")

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.runAddr.host, a.runAddr.port),
		Handler: a.router(),
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return server.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		a.log.Info("Graceful shutting down server")
		return server.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (a application) router() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			resp := struct {
				Title   string `json:"title"`
				Version string `json:"version"`
			}{
				Title:   "Накопительная система лояльности 'Гофермарт'",
				Version: "1.0",
			}

			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})
	})

	return r
}
