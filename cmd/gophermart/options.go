package main

import (
	"net/http"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
)

type option func(a *application)

func withRunAddr(host string, port int) option {
	return func(a *application) {
		a.runAddr = addr{
			host: host,
			port: port,
		}
	}
}

func withLogger(log logger.Logger) option {
	return func(a *application) {
		a.log = log
	}
}

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
