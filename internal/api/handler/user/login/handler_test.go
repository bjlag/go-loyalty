package login_test

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
	"github.com/bjlag/go-loyalty/internal/infrastructure/auth"
	"github.com/bjlag/go-loyalty/internal/infrastructure/logger/mock"
	mockRep "github.com/bjlag/go-loyalty/internal/infrastructure/repository/mock"
	"github.com/bjlag/go-loyalty/internal/model"
	ucLogin "github.com/bjlag/go-loyalty/internal/usecase/user/login"
)

func TestHandler_Handle(t *testing.T) {
	hasher := auth.NewHasher()
	jwtBuilder := auth.NewJWTBuilder("secret", 1*time.Hour)

	type args struct {
		repo func(ctrl *gomock.Controller) *mockRep.MockUserRepository
		log  func(ctrl *gomock.Controller) *mock.MockLogger
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
				repo: func(ctrl *gomock.Controller) *mockRep.MockUserRepository {
					repUserMock := mockRep.NewMockUserRepository(ctrl)
					repUserMock.EXPECT().FindByLogin(gomock.Any(), "abcd").Return(&model.User{
						GUID:     "41d2f86c-6ce5-4732-a485-6d09d7a9b3f7",
						Login:    "abcd",
						Password: "$2a$10$wEwL0jTt5ryuBRzCv56A3eq0odey9nSFrcuqqubJttyLjAw3SF2/.",
					}, nil)

					return repUserMock
				},
				log: mock.NewMockLogger,
			},
			body: `{"login": "abcd", "password": "123456"}`,
			want: want{
				status: http.StatusOK,
				err:    false,
			},
		},
		{
			name: "wrong_password",
			args: args{
				repo: func(ctrl *gomock.Controller) *mockRep.MockUserRepository {
					repUserMock := mockRep.NewMockUserRepository(ctrl)
					repUserMock.EXPECT().FindByLogin(gomock.Any(), "abcd").Return(&model.User{
						GUID:     "41d2f86c-6ce5-4732-a485-6d09d7a9b3f7",
						Login:    "abcd",
						Password: "$2a$10$wEwL0jTt5ryuBRzCv56A3eq0odey9nSFrcuqqubJttyLjAw3SF2/.",
					}, nil)

					return repUserMock
				},
				log: mock.NewMockLogger,
			},
			body: `{"login": "abcd", "password": "wrong"}`,
			want: want{
				status: http.StatusUnauthorized,
				err:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			usecase := ucLogin.NewUsecase(tt.args.repo(ctrl), hasher, jwtBuilder)
			handler := http.HandlerFunc(login.NewHandler(usecase, tt.args.log(ctrl)).Handle)

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
