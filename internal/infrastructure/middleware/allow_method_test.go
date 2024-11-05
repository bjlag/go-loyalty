package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-loyalty/internal/infrastructure/middleware"
)

func TestAllowMethod(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		wantStatusCode int
	}{
		{
			name:           "POST",
			method:         http.MethodPost,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "GET",
			method:         http.MethodGet,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "PUT",
			method:         http.MethodPut,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "PATCH",
			method:         http.MethodPatch,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "DELETE",
			method:         http.MethodDelete,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, "/", nil)

			h := middleware.AllowMethods(http.MethodPost, http.MethodGet)(http.HandlerFunc(handlerAllowMethod))
			h.ServeHTTP(w, request)

			response := w.Result()
			defer func() {
				_ = response.Body.Close()
			}()

			assert.Equal(t, tt.wantStatusCode, response.StatusCode)
		})

	}
}

func handlerAllowMethod(_ http.ResponseWriter, _ *http.Request) {}
