package create_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mockRep "github.com/bjlag/go-loyalty/internal/infrastructure/repository/mock"
	"github.com/bjlag/go-loyalty/internal/model"
	"github.com/bjlag/go-loyalty/internal/usecase/accrual/create"
)

func TestUsecase_CreateAccrual(t *testing.T) {
	var errSomeError = errors.New("some error")

	accrual := &model.Accrual{
		OrderNumber: "12345678903",
		UserGUID:    "user-123",
		Status:      model.New,
		UploadedAt:  time.Now(),
	}

	type fields struct {
		repo func(ctrl *gomock.Controller) *mockRep.MockAccrualRepo
	}

	type args struct {
		accrual *model.Accrual
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				repo: func(ctrl *gomock.Controller) *mockRep.MockAccrualRepo {
					repoMock := mockRep.NewMockAccrualRepo(ctrl)

					gomock.InOrder(
						repoMock.EXPECT().AccrualByOrderNumber(gomock.Any(), "12345678903").Return(nil, nil),
						repoMock.EXPECT().Create(gomock.Any(), accrual).Return(nil),
					)

					return repoMock
				},
			},
			args: args{
				accrual: accrual,
			},
			wantErr: assert.NoError,
		},
		{
			name: "wrong_order_number",
			fields: fields{
				repo: mockRep.NewMockAccrualRepo,
			},
			args: args{
				accrual: &model.Accrual{OrderNumber: "123"},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, create.ErrInvalidOrderNumber)
			},
		},
		{
			name: "another_user_has_already_registered_order",
			fields: fields{
				repo: func(ctrl *gomock.Controller) *mockRep.MockAccrualRepo {
					repoMock := mockRep.NewMockAccrualRepo(ctrl)

					existAccrual := &model.Accrual{
						UserGUID: "other_user",
					}

					gomock.InOrder(
						repoMock.EXPECT().AccrualByOrderNumber(gomock.Any(), "12345678903").Return(existAccrual, nil),
						repoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					)

					return repoMock
				},
			},
			args: args{
				accrual: accrual,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, create.ErrAnotherUserHasAlreadyRegisteredOrder)
			},
		},
		{
			name: "user_already_sent_this_order",
			fields: fields{
				repo: func(ctrl *gomock.Controller) *mockRep.MockAccrualRepo {
					repoMock := mockRep.NewMockAccrualRepo(ctrl)

					existAccrual := &model.Accrual{
						UserGUID: "user-123",
					}

					gomock.InOrder(
						repoMock.EXPECT().AccrualByOrderNumber(gomock.Any(), "12345678903").Return(existAccrual, nil),
						repoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					)

					return repoMock
				},
			},
			args: args{
				accrual: accrual,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, create.ErrOrderAlreadyExists)
			},
		},
		{
			name: "error_get_accrual",
			fields: fields{
				repo: func(ctrl *gomock.Controller) *mockRep.MockAccrualRepo {
					repoMock := mockRep.NewMockAccrualRepo(ctrl)

					gomock.InOrder(
						repoMock.EXPECT().AccrualByOrderNumber(gomock.Any(), "12345678903").Return(nil, errSomeError),
						repoMock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0),
					)

					return repoMock
				},
			},
			args: args{
				accrual: accrual,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, errSomeError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			u := create.NewUsecase(tt.fields.repo(ctrl))

			err := u.CreateAccrual(context.Background(), tt.args.accrual)
			if !tt.wantErr(t, err) {
				require.Fail(t, "Received unexpected error", err)
			}
		})
	}
}
