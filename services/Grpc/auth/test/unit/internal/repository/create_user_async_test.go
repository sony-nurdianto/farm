package unit_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/stretchr/testify/assert"
)

func setupBasicMocks(ctrl *gomock.Controller) (
	*mocks.MockPostgresInstance,
	*mocks.MockPostgresDatabase,
	*mocks.MockStmt,
	*mocks.MockSchemaRegisteryInstance,
	*mocks.MockClient,
	*mocks.MockAvrSerdeInstance,
	*mocks.MockKafka,
	*mocks.MockKevProducer,
) {
	mockPgDb := mocks.NewMockPostgresDatabase(ctrl)
	mockStmt := mocks.NewMockStmt(ctrl)
	mockPgI := mocks.NewMockPostgresInstance(ctrl)

	mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
	mockClient := mocks.NewMockClient(ctrl)
	mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
	mockKev := mocks.NewMockKafka(ctrl)
	mockProducer := mocks.NewMockKevProducer(ctrl)

	return mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev, mockProducer
}

func setupRepositoryMocks(
	mockPgI *mocks.MockPostgresInstance,
	mockPgDb *mocks.MockPostgresDatabase,
	mockStmt *mocks.MockStmt,
	mockSchrgs *mocks.MockSchemaRegisteryInstance,
	mockClient *mocks.MockClient,
	mockKev *mocks.MockKafka,
	mockProducer *mocks.MockKevProducer,
) {
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

	mockProducer.EXPECT().
		InitTransactions(gomock.Any()).
		Return(nil)
}

func TestCreateUserAsync_PublishAvro_ErrorNewGenericSerializer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev, mockProducer := setupBasicMocks(ctrl)

	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockKev, mockProducer)

	eventsChan := make(chan kev.Event, 1)

	mockProducer.EXPECT().
		Events().
		Return(eventsChan).
		AnyTimes()

	mockAvr.EXPECT().
		NewGenericSerializer(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(nil, errors.New("Something Wrong"))

	rp, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Something Wrong")

	close(eventsChan)
}

func TestCreateUserAsync_PublishAvro_ErrorSerializeAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev, mockProducer := setupBasicMocks(ctrl)

	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockKev, mockProducer)

	eventsChan := make(chan kev.Event, 1)

	mockProducer.EXPECT().
		Events().
		Return(eventsChan).
		AnyTimes()

	mockSerializer := mocks.NewMockAvrSerializer(ctrl)

	mockAvr.EXPECT().
		NewGenericSerializer(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(mockSerializer, nil)

	mockSerializer.EXPECT().
		Serialize(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("Failed Serialize Accounts")).
		Times(1)

	rp, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed Serialize Accounts")

	close(eventsChan)
}

func TestCreateUserAsync_PublishAvro_ErrorSerializeUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev, mockProducer := setupBasicMocks(ctrl)

	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockKev, mockProducer)

	eventsChan := make(chan kev.Event, 1)

	mockProducer.EXPECT().
		Events().
		Return(eventsChan).
		AnyTimes()

	mockSerializer := mocks.NewMockAvrSerializer(ctrl)

	mockAvr.EXPECT().
		NewGenericSerializer(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(mockSerializer, nil)

	mockSerializer.EXPECT().
		Serialize(gomock.Any(), gomock.Any()).
		Return([]byte("Accounts"), nil).
		Times(1)

	mockSerializer.EXPECT().
		Serialize(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("Failed To Serialize Users")).
		Times(1)

	rp, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed To Serialize Users")
	close(eventsChan)
}

func TestCreateUserAsync_PublishAvro_BeginTransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev, mockProducer := setupBasicMocks(ctrl)

	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockKev, mockProducer)

	mockSerializer := mocks.NewMockAvrSerializer(ctrl)

	mockAvr.EXPECT().
		NewGenericSerializer(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(mockSerializer, nil)

	mockSerializer.EXPECT().
		Serialize(gomock.Any(), gomock.Any()).
		Return([]byte("Accounts and Users"), nil).
		Times(2)

	eventsChan := make(chan kev.Event, 1)
	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()
	mockProducer.EXPECT().BeginTransaction().Return(errors.New("Error BeginTransaction")).AnyTimes()

	rp, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	// Run test
	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.EqualError(t, err, "Error BeginTransaction")

	close(eventsChan) // penting supaya goroutine exit}
}

func TestCreateUserAsync_PublishAvro_ProduceAccountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev, mockProducer := setupBasicMocks(ctrl)

	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockKev, mockProducer)

	mockSerializer := mocks.NewMockAvrSerializer(ctrl)

	mockAvr.EXPECT().
		NewGenericSerializer(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(mockSerializer, nil)

	mockSerializer.EXPECT().
		Serialize(gomock.Any(), gomock.Any()).
		Return([]byte("Accounts and Users"), nil).
		Times(2)

	eventsChan := make(chan kev.Event, 1)
	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()
	mockProducer.EXPECT().BeginTransaction().Return(nil).AnyTimes()

	// Produce for accountRecord returns error
	mockProducer.EXPECT().
		Produce(gomock.Any(), gomock.Nil()).
		Return(fmt.Errorf("Produce account error"))

	// AbortTransaction should be called after produce error
	mockProducer.EXPECT().
		AbortTransaction(gomock.Any()).
		Return(nil)

	rp, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Produce account error")

	close(eventsChan)
}

func TestCreateUserAsync_PublishAvro_ProduceUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev, mockProducer := setupBasicMocks(ctrl)

	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockKev, mockProducer)

	mockSerializer := mocks.NewMockAvrSerializer(ctrl)

	mockAvr.EXPECT().
		NewGenericSerializer(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(mockSerializer, nil)

	mockSerializer.EXPECT().
		Serialize(gomock.Any(), gomock.Any()).
		Return([]byte("Accounts and Users"), nil).
		Times(2)

	eventsChan := make(chan kev.Event, 1)
	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()
	mockProducer.EXPECT().BeginTransaction().Return(nil).AnyTimes()

	// Produce for accountRecord success
	mockProducer.EXPECT().
		Produce(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	// Produce for userRecord error
	mockProducer.EXPECT().
		Produce(gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("Produce user error")).
		Times(1)

	// AbortTransaction should be called after produce error
	mockProducer.EXPECT().
		AbortTransaction(gomock.Any()).
		Return(nil)

	rp, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Produce user error")

	close(eventsChan)
}

func TestCreateUserAsync_PublishAvro_CommitSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev, mockProducer := setupBasicMocks(ctrl)

	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockKev, mockProducer)

	mockSerializer := mocks.NewMockAvrSerializer(ctrl)

	mockAvr.EXPECT().
		NewGenericSerializer(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(mockSerializer, nil)

	mockSerializer.EXPECT().
		Serialize(gomock.Any(), gomock.Any()).
		Return([]byte("Accounts and Users"), nil).
		Times(2)

	eventsChan := make(chan kev.Event, 1)
	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()
	mockProducer.EXPECT().BeginTransaction().Return(nil).AnyTimes()

	// Produce both success
	mockProducer.EXPECT().
		Produce(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(2)

	// CommitTransaction success
	mockProducer.EXPECT().
		CommitTransaction(gomock.Any()).
		Return(nil)

	rp, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.NoError(t, err)

	close(eventsChan)
}
