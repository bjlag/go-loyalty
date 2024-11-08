//go:generate mockgen -package mock -destination mock/logger_mock.go github.com/bjlag/go-loyalty/internal/infrastructure/logger Logger

package middleware

import (
	"net/http"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
)

func LogRequest(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			dw := newResponseDataWriter(w)

			start := time.Now()
			next.ServeHTTP(dw, r)
			duration := time.Since(start)

			log.
				WithField("request_id", chiMiddleware.GetReqID(r.Context())).
				WithField("method", r.Method).
				WithField("uri", r.URL.Path).
				WithField("status", dw.data.status).
				WithField("duration", duration).
				WithField("size", dw.data.size).
				Info("Got request")
		}

		return http.HandlerFunc(fn)
	}
}

type responseData struct {
	status int
	size   int
}

type responseDataWriter struct {
	http.ResponseWriter

	data *responseData
}

func newResponseDataWriter(w http.ResponseWriter) *responseDataWriter {
	return &responseDataWriter{
		ResponseWriter: w,
		data: &responseData{
			status: http.StatusOK,
			size:   0,
		},
	}
}

func (w *responseDataWriter) Write(buf []byte) (int, error) {
	size, err := w.ResponseWriter.Write(buf)
	w.data.size += size
	return size, err
}

func (w *responseDataWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.data.status = status
}
