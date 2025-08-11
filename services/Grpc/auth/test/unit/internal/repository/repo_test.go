package unit_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
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

		mockProducer := mocks.NewMockKevProducer(ctrl)

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

		mockSchrgs.EXPECT().
			NewClient(gomock.Any()).
			Return(mockClient, nil)

		mockKev.EXPECT().
			NewProducer(gomock.Any()).
			Return(mockProducer, nil)

		eventsChan := make(chan kev.Event, 1) // buffered supaya gak blocking

		mockProducer.EXPECT().
			Events().
			Return(eventsChan).
			AnyTimes()

		mockProducer.EXPECT().
			InitTransactions(gomock.Any()).
			Return(nil)

		_, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
		assert.NoError(t, err)
		close(eventsChan)
	})

	t.Run("NewPostgreRepo Error OpenPostgres", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)

		mockKev := mocks.NewMockKafka(ctrl)

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

		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
		mockKev := mocks.NewMockKafka(ctrl)

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
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
		mockKev := mocks.NewMockKafka(ctrl)

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

	t.Run("NewPostgreRepo Error PrepareStmt GetUserEmail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockPgDb := mocks.NewMockPostgresDatabase(ctrl)
		mockStmt := mocks.NewMockStmt(ctrl)

		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
		mockKev := mocks.NewMockKafka(ctrl)

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

	t.Run("NewPostgreRepo Error NewSchemaRegistery", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockPgDb := mocks.NewMockPostgresDatabase(ctrl)
		mockStmt := mocks.NewMockStmt(ctrl)

		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
		mockKev := mocks.NewMockKafka(ctrl)

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

		mockSchrgs.EXPECT().
			NewClient(gomock.Any()).
			Return(nil, errors.New("Failed To Create NewClient"))

		_, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
		assert.Error(t, err)
		assert.EqualError(t, err, "Failed To Create NewClient")
	})

	t.Run("NewPostgreRepo Error Producer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockPgDb := mocks.NewMockPostgresDatabase(ctrl)
		mockStmt := mocks.NewMockStmt(ctrl)

		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
		mockKev := mocks.NewMockKafka(ctrl)

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

		mockSchrgs.EXPECT().
			NewClient(gomock.Any()).
			Return(nil, errors.New("Failed Create Producer"))

		_, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
		assert.Error(t, err)
		assert.EqualError(t, err, "Failed Create Producer")
	})

	t.Run("NewPostgreRepo Error InitTransactions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI := mocks.NewMockPostgresInstance(ctrl)
		mockPgDb := mocks.NewMockPostgresDatabase(ctrl)
		mockStmt := mocks.NewMockStmt(ctrl)

		mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockClient := mocks.NewMockClient(ctrl)
		mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
		mockKev := mocks.NewMockKafka(ctrl)

		mockProducer := mocks.NewMockKevProducer(ctrl)

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

		mockSchrgs.EXPECT().
			NewClient(gomock.Any()).
			Return(mockClient, nil)

		mockKev.EXPECT().
			NewProducer(gomock.Any()).
			Return(mockProducer, nil)

		eventsChan := make(chan kev.Event, 1) // buffered supaya gak blocking

		mockProducer.EXPECT().
			Events().
			Return(eventsChan).
			AnyTimes()

		mockProducer.EXPECT().
			InitTransactions(gomock.Any()).
			Return(errors.New("Error Init Transactions")).AnyTimes()

		_, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
		assert.Error(t, err)
		assert.EqualError(t, err, "init transactions failed after 5 attempts: Error Init Transactions")
		close(eventsChan)
	})
}
