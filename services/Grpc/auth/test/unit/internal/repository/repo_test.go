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
	t.Run("NewPostgreRepo Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockPgDb := mocks.NewMockPostgresDatabase(ctrl)
		mockStmt := mocks.NewMockStmt(ctrl)

		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockClient := mocks.NewMockClient(ctrl)
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
		mockKev := mocks.NewMockKafka(ctrl)

		mockSchrgs.EXPECT().
			NewClient(gomock.Any()).
			Return(mockClient, nil)

		mockPgI.EXPECT().
			Open(gomock.Any(), gomock.Any()).
			Return(mockPgDb, nil)

		mockPgDb.EXPECT().
			SetMaxOpenConns(gomock.Any()).
			Return()

		mockPgDb.EXPECT().
			SetMaxIdleConns(gomock.Any()).
			Return()

		mockPgDb.EXPECT().
			SetConnMaxLifetime(gomock.Any()).
			Return()

		mockPgDb.EXPECT().
			PingContext(gomock.Any()).
			Return(nil)

		mockPgDb.EXPECT().
			Prepare(gomock.Any()).
			Return(mockStmt, nil).Times(2)

		_, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
		assert.NoError(t, err)
	})

	t.Run("NewPostgreRepo Error NewSchemaRegistery", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
		mockKev := mocks.NewMockKafka(ctrl)

		mockSchrgs.EXPECT().
			NewClient(gomock.Any()).
			Return(nil, errors.New("Failed To Create NewClient"))

		res, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("NewPostgreRepo Error OpenPostgres", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockClient := mocks.NewMockClient(ctrl)
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)

		mockKev := mocks.NewMockKafka(ctrl)

		mockSchrgs.EXPECT().
			NewClient(gomock.Any()).
			Return(mockClient, nil)

		mockPgI.EXPECT().
			Open(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("Failed To Open Postgres"))

		_, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
		assert.Error(t, err)
	})

	t.Run("NewPostgreRepo Error PrepareStmt CreateUser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockPgDb := mocks.NewMockPostgresDatabase(ctrl)
		// mockStmt := mocks.NewMockStmt(ctrl)

		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockClient := mocks.NewMockClient(ctrl)
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
		mockKev := mocks.NewMockKafka(ctrl)

		mockSchrgs.EXPECT().
			NewClient(gomock.Any()).
			Return(mockClient, nil)

		mockPgI.EXPECT().
			Open(gomock.Any(), gomock.Any()).
			Return(mockPgDb, nil)

		mockPgDb.EXPECT().
			SetMaxOpenConns(gomock.Any()).
			Return()

		mockPgDb.EXPECT().
			SetMaxIdleConns(gomock.Any()).
			Return()

		mockPgDb.EXPECT().
			SetConnMaxLifetime(gomock.Any()).
			Return()

		mockPgDb.EXPECT().
			PingContext(gomock.Any()).
			Return(nil)

		mockPgDb.EXPECT().
			Prepare(gomock.Any()).
			Return(nil, errors.New("Failed To Create STMT"))

		_, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
		assert.Error(t, err)
		assert.EqualError(t, err, "Failed To Create STMT")
	})

	t.Run("NewPostgreRepo Error PrepareStmt GetUserEmail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockPgDb := mocks.NewMockPostgresDatabase(ctrl)
		mockStmt := mocks.NewMockStmt(ctrl)

		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockClient := mocks.NewMockClient(ctrl)
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
		mockKev := mocks.NewMockKafka(ctrl)

		mockSchrgs.EXPECT().
			NewClient(gomock.Any()).
			Return(mockClient, nil)

		mockPgI.EXPECT().
			Open(gomock.Any(), gomock.Any()).
			Return(mockPgDb, nil)

		mockPgDb.EXPECT().
			SetMaxOpenConns(gomock.Any()).
			Return()

		mockPgDb.EXPECT().
			SetMaxIdleConns(gomock.Any()).
			Return()

		mockPgDb.EXPECT().
			SetConnMaxLifetime(gomock.Any()).
			Return()

		mockPgDb.EXPECT().
			PingContext(gomock.Any()).
			Return(nil)

		mockPgDb.EXPECT().
			Prepare(gomock.Any()).
			Return(mockStmt, nil)

		mockPgDb.EXPECT().
			Prepare(gomock.Any()).
			Return(nil, errors.New("Failed to Create STMT GetUSerEmail"))

		_, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
		assert.Error(t, err)
		assert.EqualError(t, err, "Failed to Create STMT GetUSerEmail")
	})
}
