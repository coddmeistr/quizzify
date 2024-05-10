package testsservice

import (
	"context"
	"errors"
	"github.com/coddmeistr/quizzify/backend/tests/internal/config"
	"github.com/coddmeistr/quizzify/backend/tests/internal/domain"
	"github.com/coddmeistr/quizzify/backend/tests/internal/helpers/user"
	"github.com/coddmeistr/quizzify/backend/tests/internal/service/tests/mocks"
	"github.com/coddmeistr/quizzify/backend/tests/internal/storage"
	p "github.com/coddmeistr/quizzify/backend/tests/pkg/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestService_CreateTest(t *testing.T) {
	testType := reflect.TypeOf(domain.Test{}).String()

	type args struct {
		test domain.Test
	}
	tc := []struct {
		name      string
		args      args
		wantError bool
		err       error
		mockF     func(st *mocks.Storage, val *mocks.Validator)
	}{
		{
			name: "ok",
			args: args{
				test: domain.Test{},
			},
			wantError: false,
			err:       nil,
			mockF: func(st *mocks.Storage, val *mocks.Validator) {
				val.On("ValidateTest", mock.AnythingOfType(testType)).Return(nil).Once()
				st.On("CreateTest", context.Background(), mock.AnythingOfType(testType)).Return(nil).Once()
			},
		},
		{
			name: "failed validation",
			args: args{
				test: domain.Test{},
			},
			wantError: true,
			err:       errors.New("failed test validation"),
			mockF: func(st *mocks.Storage, val *mocks.Validator) {
				val.On("ValidateTest", mock.AnythingOfType(testType)).Return(errors.New("failed test validation")).Once()
			},
		},
		{
			name: "error in storage",
			args: args{
				test: domain.Test{},
			},
			wantError: true,
			err:       errors.New("some error occured"),
			mockF: func(st *mocks.Storage, val *mocks.Validator) {
				val.On("ValidateTest", mock.AnythingOfType(testType)).Return(nil).Once()
				st.On("CreateTest", context.Background(), mock.AnythingOfType(testType)).Return(errors.New("some error occured")).Once()
			},
		},
	}

	for _, tt := range tc {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockVal := mocks.NewValidator(t)
			mockSt := mocks.NewStorage(t)
			tt.mockF(mockSt, mockVal)

			s := New(zap.NewExample(), &config.Config{}, mockSt, mockVal)
			_, got := s.CreateTest(context.Background(), tt.args.test)

			if tt.wantError {
				assert.Containsf(t, got.Error(), tt.err.Error(), "expected error containing %q, got %s", tt.err.Error(), got.Error())
			} else {
				assert.NoError(t, got)
			}
		})
	}

}

func TestService_DeleteTest(t *testing.T) {
	validTestId := "623452gsgsgf"

	type args struct {
		ctx    context.Context
		testID string
	}
	tc := []struct {
		name      string
		args      args
		wantError bool
		err       error
		mockF     func(st *mocks.Storage)
	}{
		{
			name: "ok",
			args: args{
				ctx: context.WithValue(context.Background(), user.AuthInfoKey, user.Info{
					ID: 1,
				}),
				testID: validTestId,
			},
			wantError: false,
			err:       nil,
			mockF: func(st *mocks.Storage) {
				gotTest := &domain.Test{
					ID:     &validTestId,
					UserID: p.Int(1),
				}
				st.On("GetTestByID", mock.Anything, validTestId, false).Return(gotTest, nil).Once()
				st.On("DeleteTest", mock.Anything, validTestId).Return(nil).Once()
			},
		},
		{
			name: "ok, but accesed with admin rights",
			args: args{
				ctx: context.WithValue(context.Background(), user.AuthInfoKey, user.Info{
					ID:          1,
					Permissions: []int{user.Admin},
				}),
				testID: validTestId,
			},
			wantError: false,
			err:       nil,
			mockF: func(st *mocks.Storage) {
				gotTest := &domain.Test{
					ID:     &validTestId,
					UserID: p.Int(2),
				}
				st.On("GetTestByID", mock.Anything, validTestId, false).Return(gotTest, nil).Once()
				st.On("DeleteTest", mock.Anything, validTestId).Return(nil).Once()
			},
		},
		{
			name: "no rights to perform",
			args: args{
				ctx: context.WithValue(context.Background(), user.AuthInfoKey, user.Info{
					ID: 1,
				}),
				testID: validTestId,
			},
			wantError: true,
			err:       errors.New("no rights to perform"),
			mockF: func(st *mocks.Storage) {
				gotTest := &domain.Test{
					ID:     &validTestId,
					UserID: p.Int(2),
				}
				st.On("GetTestByID", mock.Anything, validTestId, false).Return(gotTest, nil).Once()
			},
		},
		{
			name: "not found test to delete",
			args: args{
				ctx: context.WithValue(context.Background(), user.AuthInfoKey, user.Info{
					ID: 1,
				}),
				testID: validTestId,
			},
			wantError: true,
			err:       errors.New("not found"),
			mockF: func(st *mocks.Storage) {
				st.On("GetTestByID", mock.Anything, validTestId, false).Return(nil, storage.ErrNotFound).Once()
			},
		},
		{
			name: "error while getting test to delete",
			args: args{
				ctx: context.WithValue(context.Background(), user.AuthInfoKey, user.Info{
					ID: 1,
				}),
				testID: validTestId,
			},
			wantError: true,
			err:       errors.New("some error while getting test"),
			mockF: func(st *mocks.Storage) {
				st.On("GetTestByID", mock.Anything, validTestId, false).Return(nil, errors.New("some error while getting test")).Once()
			},
		},
		{
			name: "error while deleting test",
			args: args{
				ctx: context.WithValue(context.Background(), user.AuthInfoKey, user.Info{
					ID: 1,
				}),
				testID: validTestId,
			},
			wantError: true,
			err:       errors.New("some error while deleting test"),
			mockF: func(st *mocks.Storage) {
				gotTest := &domain.Test{
					ID:     &validTestId,
					UserID: p.Int(1),
				}
				st.On("GetTestByID", mock.Anything, validTestId, false).Return(gotTest, nil).Once()
				st.On("DeleteTest", mock.Anything, validTestId).Return(errors.New("some error while deleting test")).Once()
			},
		},
		{
			name: "not found test when deleting",
			args: args{
				ctx: context.WithValue(context.Background(), user.AuthInfoKey, user.Info{
					ID: 1,
				}),
				testID: validTestId,
			},
			wantError: true,
			err:       storage.ErrNotFound,
			mockF: func(st *mocks.Storage) {
				gotTest := &domain.Test{
					ID:     &validTestId,
					UserID: p.Int(1),
				}
				st.On("GetTestByID", mock.Anything, validTestId, false).Return(gotTest, nil).Once()
				st.On("DeleteTest", mock.Anything, validTestId).Return(storage.ErrNotFound).Once()
			},
		},
	}

	for _, tt := range tc {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockVal := mocks.NewValidator(t)
			mockSt := mocks.NewStorage(t)
			tt.mockF(mockSt)

			s := New(zap.NewExample(), &config.Config{}, mockSt, mockVal)
			got := s.DeleteTest(tt.args.ctx, tt.args.testID)

			if tt.wantError {
				assert.Containsf(t, got.Error(), tt.err.Error(), "expected error containing %q, got %s", tt.err.Error(), got.Error())
			} else {
				assert.NoError(t, got)
			}
		})
	}

}
