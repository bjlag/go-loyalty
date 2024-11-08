package middleware_test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger/mock"
	"github.com/bjlag/go-loyalty/internal/infrastructure/middleware"
)

func TestGzip(t *testing.T) {
	type body struct {
		Value int `json:"value"`
	}

	compressBody := func(t *testing.T, body body) []byte {
		var b bytes.Buffer
		w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
		require.NoError(t, err)

		marshaledBody, err := json.Marshal(body)
		require.NoError(t, err)

		_, err = w.Write(marshaledBody)
		require.NoError(t, err)
		err = w.Close()
		require.NoError(t, err)

		return b.Bytes()
	}

	decompressBody := func(t *testing.T, r io.Reader) body {
		zr, err := gzip.NewReader(r)
		require.NoError(t, err)

		var zb bytes.Buffer
		_, err = zb.ReadFrom(zr)
		require.NoError(t, err)

		var respBody body
		err = json.Unmarshal(zb.Bytes(), &respBody)
		require.NoError(t, err)

		return respBody
	}

	t.Run("compress", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		var reqBody body
		reqBody.Value = 1

		compressedBody := compressBody(t, reqBody)

		w := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/url", bytes.NewReader(compressedBody))
		request.Header.Set("Content-Encoding", "gzip")
		request.Header.Set("Accept-Encoding", "gzip")

		h := middleware.Gzip(mock.NewMockLogger(ctrl))(http.HandlerFunc(handlerGzip))
		h.ServeHTTP(w, request)

		response := w.Result()
		defer func() {
			_ = response.Body.Close()
		}()

		respBody := decompressBody(t, response.Body)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, 2, respBody.Value)
	})

	t.Run("without_compress", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		var reqBody body
		reqBody.Value = 1

		marshaledBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/url", bytes.NewReader(marshaledBody))

		h := middleware.Gzip(mock.NewMockLogger(ctrl))(http.HandlerFunc(handlerGzip))
		h.ServeHTTP(w, request)

		response := w.Result()
		defer func() {
			_ = response.Body.Close()
		}()

		var respBody body

		err = json.NewDecoder(response.Body).Decode(&respBody)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, 2, respBody.Value)
	})
}

func handlerGzip(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Value int `json:"value"`
	}

	_ = json.NewDecoder(r.Body).Decode(&body)
	body.Value += 1
	_ = json.NewEncoder(w).Encode(body)

	w.WriteHeader(http.StatusOK)
}
