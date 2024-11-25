package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/bjlag/go-loyalty/internal/infrastructure/logger/mock"
	"github.com/bjlag/go-loyalty/internal/infrastructure/middleware"
)

func TestLogRequest(t *testing.T) {
	type args struct {
		log func(ctrl *gomock.Controller) *mock.MockLogger
	}

	tests := []struct {
		name string
		args args
		want func(next http.Handler) http.Handler
	}{
		{
			name: "test",
			args: args{
				log: func(ctrl *gomock.Controller) *mock.MockLogger {
					mockLog := mock.NewMockLogger(ctrl)
					mockLog.EXPECT().WithField("request_id", gomock.Any()).Return(mockLog).AnyTimes()
					mockLog.EXPECT().WithField("method", "POST").Return(mockLog).AnyTimes()
					mockLog.EXPECT().WithField("uri", "/url").Return(mockLog).AnyTimes()
					mockLog.EXPECT().WithField("status", 200).Return(mockLog).AnyTimes()
					mockLog.EXPECT().WithField("duration", gomock.Any()).Return(mockLog).AnyTimes()
					mockLog.EXPECT().WithField("size", 7).Return(mockLog).AnyTimes()
					mockLog.EXPECT().Info("Got request").AnyTimes()

					return mockLog
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			w := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/url", nil)

			h := middleware.LogRequest(tt.args.log(ctrl))(http.HandlerFunc(handlerLogRequest))
			h.ServeHTTP(w, request)
		})
	}
}

func handlerLogRequest(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("handler"))
	w.WriteHeader(http.StatusOK)
}
