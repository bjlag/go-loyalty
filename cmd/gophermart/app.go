package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/infrastructure/middleware"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

const appVersion = "1.0.0"

var errNoLogger = errors.New("no logger provided (use 'withLogger' option)")

type runAddr struct {
	host string
	port int
}

type apiHandler struct {
	method      string
	path        string
	handler     http.HandlerFunc
	middlewares []func(next http.Handler) http.Handler
}

type application struct {
	runAddr     runAddr
	log         logger.Logger
	apiHandlers []apiHandler
}

func newApp(opts ...option) *application {
	a := &application{}
	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a application) run(ctx context.Context) error {
	if a.log == nil {
		return errNoLogger
	}

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

	r.Use(
		chiMiddleware.RequestID,
		middleware.LogRequest(a.log),
		middleware.Gzip(a.log),
	)

	for _, h := range a.apiHandlers {
		r.With(h.middlewares...).Method(h.method, h.path, h.handler)
	}

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp := struct {
			Title   string `json:"title"`
			Version string `json:"version"`
		}{
			Title:   "Накопительная система лояльности 'Гофермарт'",
			Version: appVersion,
		}

		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})

	return r
}
