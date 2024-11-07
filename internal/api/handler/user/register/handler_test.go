package register_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bjlag/go-loyalty/internal/api/handler/user/login"
	"github.com/bjlag/go-loyalty/internal/api/handler/user/register"
	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	mockAuth "github.com/bjlag/go-loyalty/internal/infrastructure/auth/mock"
	mockGuid "github.com/bjlag/go-loyalty/internal/infrastructure/guid/mock"
	mockLog "github.com/bjlag/go-loyalty/internal/infrastructure/logger/mock"
	mockRep "github.com/bjlag/go-loyalty/internal/infrastructure/repository/mock"
	"github.com/bjlag/go-loyalty/internal/model"
	ucRegister "github.com/bjlag/go-loyalty/internal/usecase/user/register"
)

func TestHandler_Handle(t *testing.T) {
	jwtBuilder := auth.NewJWTBuilder("secret", 1*time.Hour)

	type args struct {
		rep       func(ctrl *gomock.Controller) *mockRep.MockUserRepository
		hasher    func(ctrl *gomock.Controller) *mockAuth.MockIHasher
		generator func(ctrl *gomock.Controller) *mockGuid.MockIGenerator
		log       func(ctrl *gomock.Controller) *mockLog.MockLogger
	}

	type want struct {
		status int
		err    bool
	}

	tests := []struct {
		name string
		args args
		body string
		want want
	}{
		{
			name: "success",
			args: args{
				rep: func(ctrl *gomock.Controller) *mockRep.MockUserRepository {
					repUserMock := mockRep.NewMockUserRepository(ctrl)

					user := &model.User{
						GUID:     "41d2f86c-6ce5-4732-a485-6d09d7a9b3f7",
						Login:    "abcd",
						Password: "$2a$10$wEwL0jTt5ryuBRzCv56A3eq0odey9nSFrcuqqubJttyLjAw3SF2/.",
					}

					gomock.InOrder(
						repUserMock.EXPECT().FindByLogin(gomock.Any(), "abcd").Return(nil, nil),
						repUserMock.EXPECT().Insert(gomock.Any(), user).Return(nil),
					)

					return repUserMock
				},
				generator: func(ctrl *gomock.Controller) *mockGuid.MockIGenerator {
					genMock := mockGuid.NewMockIGenerator(ctrl)
					genMock.EXPECT().Generate().Return("41d2f86c-6ce5-4732-a485-6d09d7a9b3f7")
					return genMock
				},
				hasher: func(ctrl *gomock.Controller) *mockAuth.MockIHasher {
					hasherMock := mockAuth.NewMockIHasher(ctrl)
					hasherMock.EXPECT().HashPassword("123456").Return("$2a$10$wEwL0jTt5ryuBRzCv56A3eq0odey9nSFrcuqqubJttyLjAw3SF2/.", nil)
					return hasherMock
				},
				log: mockLog.NewMockLogger,
			},
			body: `{"login": "abcd", "password": "123456"}`,
			want: want{
				status: http.StatusOK,
				err:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			usecase := ucRegister.NewUsecase(tt.args.rep(ctrl), tt.args.generator(ctrl), tt.args.hasher(ctrl), jwtBuilder)
			handler := http.HandlerFunc(register.NewHandler(usecase, tt.args.log(ctrl)).Handle)

			srv := httptest.NewServer(handler)
			defer srv.Close()

			req := resty.New().R()
			req.Method = http.MethodPost
			req.URL = srv.URL
			req.SetHeader("Content-Type", "application/json")
			req.SetBody(tt.body)

			resp, err := req.Send()
			require.NoError(t, err, "error making HTTP request")

			if !tt.want.err {
				var respUnmarshalled login.Response
				err = json.Unmarshal(resp.Body(), &respUnmarshalled)
				require.NoError(t, err, "error unmarshaling HTTP response body")

				assert.NotEmpty(t, respUnmarshalled.Token, "token is empty")
				assert.Equal(t, resp.Header().Get("Content-Type"), "application/json")
				assert.Equal(t, resp.Header().Get("Authorization"), fmt.Sprintf("Bearer %s", respUnmarshalled.Token))
			}

			assert.Equal(t, tt.want.status, resp.StatusCode(), "unexpected status code")
		})
	}
}
