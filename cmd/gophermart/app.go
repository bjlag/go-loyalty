package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
	"github.com/bjlag/go-loyalty/internal/infrastructure/middleware"
)

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

type option func(a *application)

func withAPIHandler(method, path string, handler http.HandlerFunc, middlewares ...func(next http.Handler) http.Handler) option {
	return func(a *application) {
		a.apiHandlers = append(a.apiHandlers, apiHandler{
			method:      method,
			path:        path,
			handler:     handler,
			middlewares: middlewares,
		})
	}
}

func newApp(runAddr runAddr, log logger.Logger, opts ...option) *application {
	// todo если лог не передан, то делаем по умолчанию log.
	// todo run адрес передавать передавать через опции

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

	r.Use(
		middleware.LogRequest(a.log),
		middleware.Gzip(a.log),
	)

	for _, h := range a.apiHandlers {
		r.With(h.middlewares...).Method(h.method, h.path, h.handler)
	}

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
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		})
	})

	return r
}
