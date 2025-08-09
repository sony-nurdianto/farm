package unit_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
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
) {
	// Hanya membuat mock objects tanpa expectation
	mockPgI := mocks.NewMockPostgresInstance(ctrl)
	mockPgDb := mocks.NewMockPostgresDatabase(ctrl)
	mockStmt := mocks.NewMockStmt(ctrl)
	mockSchrgs := mocks.NewMockSchemaRegisteryInstance(ctrl)
	mockClient := mocks.NewMockClient(ctrl)
	mockAvr := mocks.NewMockAvrSerdeInstance(ctrl)
	mockKev := mocks.NewMockKafka(ctrl)

	return mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev
}

func setupRepositoryMocks(
	mockPgI *mocks.MockPostgresInstance,
	mockPgDb *mocks.MockPostgresDatabase,
	mockStmt *mocks.MockStmt,
	mockSchrgs *mocks.MockSchemaRegisteryInstance,
	mockClient *mocks.MockClient,
) {
	// Setup expectations untuk repository initialization
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
}

func TestCreateUserAsync_ErrorGetLatestSchemaMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)

	// Setup repository mocks
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata("insert-account--value").
		Return(schrgs.SchemaMetadata{}, errors.New("Error GetLatestSchemaMetadata")).
		Times(1)

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.EqualError(t, err, "Error GetLatestSchemaMetadata")
}

func TestCreateUserAsync_ErrorSchemaNotFoundCreateAvroSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)

	// Setup repository mocks
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata("insert-account--value").
		Return(schrgs.SchemaMetadata{}, errors.New("40401")).
		Times(1)

	mockClient.EXPECT().
		Register("insert-account--value", gomock.Any(), false).
		Return(0, errors.New("Failed Register Schema")).
		Times(1)

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed Register Schema")
}

func TestCreateUserAsync_ErrorOnUserSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)

	// Setup repository mocks
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata("insert-account--value").
		Return(schrgs.SchemaMetadata{ID: 1}, nil).
		Times(1)

	mockClient.EXPECT().
		GetLatestSchemaMetadata("insert-user--value").
		Return(schrgs.SchemaMetadata{}, errors.New("40401")).
		Times(1)

	mockClient.EXPECT().
		Register("insert-user--value", gomock.Any(), false).
		Return(0, errors.New("Failed Register User Schema")).
		Times(1)

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed Register User Schema")
}

func TestCreateUserAsync_PublishAvro_ErrorNewGenericSerializer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)

	// Setup repository mocks
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata(gomock.Any()).
		Return(schrgs.SchemaMetadata{}, nil).
		Times(2)

	mockAvr.EXPECT().
		NewGenericSerializer(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(nil, errors.New("Something Wrong"))

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Something Wrong")
}

func TestCreateUserAsync_PublishAvro_ErrorSerializeAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)

	// Setup repository mocks
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata(gomock.Any()).
		Return(schrgs.SchemaMetadata{}, nil).
		Times(2)

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

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed Serialize Accounts")
}

func TestCreateUserAsync_PublishAvro_ErrorSerializeUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)

	// Setup repository mocks
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata(gomock.Any()).
		Return(schrgs.SchemaMetadata{}, nil).
		Times(2)

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

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed To Serialize Users")
}

func TestCreateUserAsync_PublishAvro_ProducerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)

	// Setup repository mocks
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata(gomock.Any()).
		Return(schrgs.SchemaMetadata{ID: 1}, nil).
		Times(2)

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

	mockKev.EXPECT().
		NewProducer(gomock.Any()).
		Return(nil, errors.New("Error Create Producer"))

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to create kafka producer: Error Create Producer")
}

func TestCreateUserAsync_PublishAvro_InitTransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)

	// Setup repository mocks - ISOLATED untuk test ini
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata(gomock.Any()).
		Return(schrgs.SchemaMetadata{ID: 1}, nil).
		Times(2)

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

	// Buat producer mock baru untuk test ini
	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKev.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil).AnyTimes()

	// TODO:Need To Add AbstractionType
	eventsChan := make(chan kafka.Event, 1) // buffered supaya gak blocking

	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()
	mockProducer.EXPECT().InitTransactions(gomock.Any()).Return(errors.New("Error InitTransactions")).AnyTimes()
	// mockProducer.EXPECT().BeginTransaction().Return(errors.New("Error BeginTransaction")).AnyTimes()

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	// Run test
	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.EqualError(t, err, "Error InitTransactions")

	close(eventsChan) // penting supaya goroutine exit
}

func TestCreateUserAsync_PublishAvro_BeginTransactionError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)

	// Setup repository mocks - ISOLATED untuk test ini
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata(gomock.Any()).
		Return(schrgs.SchemaMetadata{ID: 1}, nil).
		Times(2)

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

	// Buat producer mock baru untuk test ini
	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKev.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil).AnyTimes()

	// TODO:Need To Add AbstractionType
	eventsChan := make(chan kafka.Event, 1) // buffered supaya gak blocking

	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()
	mockProducer.EXPECT().InitTransactions(gomock.Any()).Return(nil).AnyTimes()
	mockProducer.EXPECT().BeginTransaction().Return(errors.New("Error BeginTransaction")).AnyTimes()

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
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

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata(gomock.Any()).
		Return(schrgs.SchemaMetadata{ID: 1}, nil).
		Times(2)

	mockSerializer := mocks.NewMockAvrSerializer(ctrl)
	mockAvr.EXPECT().
		NewGenericSerializer(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockSerializer, nil)
	mockSerializer.EXPECT().
		Serialize(gomock.Any(), gomock.Any()).
		Return([]byte("payload"), nil).
		Times(2)

	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKev.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil).AnyTimes()

	// TODO:Need To Add AbstractionType
	eventsChan := make(chan kafka.Event, 1)
	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()

	mockProducer.EXPECT().InitTransactions(gomock.Any()).Return(nil).AnyTimes()
	mockProducer.EXPECT().BeginTransaction().Return(nil).AnyTimes()

	// Produce for accountRecord returns error
	mockProducer.EXPECT().
		Produce(gomock.Any(), gomock.Nil()).
		Return(fmt.Errorf("Produce account error"))

	// AbortTransaction should be called after produce error
	mockProducer.EXPECT().
		AbortTransaction(gomock.Any()).
		Return(nil)

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Produce account error")

	close(eventsChan)
}

func TestCreateUserAsync_PublishAvro_ProduceUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata(gomock.Any()).
		Return(schrgs.SchemaMetadata{ID: 1}, nil).
		Times(2)

	mockSerializer := mocks.NewMockAvrSerializer(ctrl)
	mockAvr.EXPECT().
		NewGenericSerializer(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockSerializer, nil)
	mockSerializer.EXPECT().
		Serialize(gomock.Any(), gomock.Any()).
		Return([]byte("payload"), nil).
		Times(2)

	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKev.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil).AnyTimes()

	// TODO:Need To Add AbstractionType
	eventsChan := make(chan kafka.Event, 1)
	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()

	mockProducer.EXPECT().InitTransactions(gomock.Any()).Return(nil).AnyTimes()
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

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Produce user error")

	close(eventsChan)
}

func TestCreateUserAsync_PublishAvro_CommitSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	mockClient.EXPECT().
		GetLatestSchemaMetadata(gomock.Any()).
		Return(schrgs.SchemaMetadata{ID: 1}, nil).
		Times(2)

	mockSerializer := mocks.NewMockAvrSerializer(ctrl)
	mockAvr.EXPECT().
		NewGenericSerializer(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockSerializer, nil)
	mockSerializer.EXPECT().
		Serialize(gomock.Any(), gomock.Any()).
		Return([]byte("payload"), nil).
		Times(2)

	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKev.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil).AnyTimes()

	// TODO:Need To Add AbstractionType
	eventsChan := make(chan kafka.Event, 1)
	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()

	mockProducer.EXPECT().InitTransactions(gomock.Any()).Return(nil).AnyTimes()
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

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.NoError(t, err)

	close(eventsChan)
}

// func TestCreateUserAsync_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
//
// 	mockPgI, _, _, mockSchrgs, mockClient := setupCommonMocks(ctrl)
//
// 	// Mock untuk account schema - berhasil (sudah ada)
// 	mockClient.EXPECT().
// 		GetLatestSchemaMetadata("insert-account--value").
// 		Return(schemaregistry.SchemaMetadata{ID: 1}, nil).
// 		Times(1)
//
// 	// Mock untuk user schema - berhasil (sudah ada)
// 	mockClient.EXPECT().
// 		GetLatestSchemaMetadata("insert-user--value").
// 		Return(schemaregistry.SchemaMetadata{ID: 2}, nil).
// 		Times(1)
//
// 	// TODO: Mock publishAvro method juga kalau diperlukan
// 	// mockClient.EXPECT().
// 	//     Produce(gomock.Any(), gomock.Any()).
// 	//     Return(nil).Times(2)
//
// 	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI)
// 	assert.NoError(t, err)
//
// 	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
// 	// assert.NoError(t, err)  // Uncomment setelah mock publishAvro
// 	assert.Error(t, err) // Sementara expect error karena publishAvro belum di-mock
// }
