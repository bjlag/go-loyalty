package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger"
)

func Gzip(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w

			if isRequestCompressed(r) {
				zr, err := newGzipReader(r.Body)
				if err != nil {
					log.WithError(err).Error("Error creating gzip reader")
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				r.Body = zr
			}

			if isClientSupportCompress(r) {
				zw, err := newGzipWriter(w)
				if err != nil {
					log.WithError(err).Error("Error creating gzip writer")
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				defer func() {
					err = zw.Close()
					if err != nil {
						log.WithError(err).Error("Failed to close gzip writer")
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					}
				}()

				zw.Header().Set("Content-Encoding", "gzip")

				ow = zw
			}

			next.ServeHTTP(ow, r)
		}

		return http.HandlerFunc(fn)
	}
}

func isRequestCompressed(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
}

func isClientSupportCompress(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

type gzipReader struct {
	zr io.ReadCloser
}

func newGzipReader(r io.Reader) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		zr: zr,
	}, nil
}

func (r *gzipReader) Read(p []byte) (int, error) {
	return r.zr.Read(p)
}

func (r *gzipReader) Close() error {
	return r.zr.Close()
}

type gzipWriter struct {
	http.ResponseWriter

	zw *gzip.Writer
}

func newGzipWriter(w http.ResponseWriter) (*gzipWriter, error) {
	zw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	return &gzipWriter{
		ResponseWriter: w,
		zw:             zw,
	}, nil
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.zw.Write(b)
}

func (w *gzipWriter) Close() error {
	return w.zw.Close()
}
