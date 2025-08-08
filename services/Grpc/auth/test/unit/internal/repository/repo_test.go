package unit_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRepo(t *testing.T) {
	t.Run("NewPostgreRepo Error NewSchemaRegistery", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)

		mockSchrgs.EXPECT().
			NewClient(gomock.Any()).
			Return(nil, errors.New("Failed To Create NewClient"))

		res, err := repository.NewPostgresRepo(mockSchrgs, mockPgI)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("NewPostgreRepo Error OpenPostgres", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockClient := mocks.NewMockClient(ctrl)

		mockSchrgs.EXPECT().
			NewClient(gomock.Any()).
			Return(mockClient, nil)

		mockPgI.EXPECT().
			Open(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("Failed To Open Postgres"))

		_, err := repository.NewPostgresRepo(mockSchrgs, mockPgI)
		assert.Error(t, err)
	})
}
