package unit_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
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

	eventsChan := make(chan kev.Event, 1) // buffered supaya gak blocking

	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()
	mockProducer.EXPECT().InitTransactions(gomock.Any()).Return(errors.New("Error InitTransactions")).AnyTimes()
	// mockProducer.EXPECT().BeginTransaction().Return(errors.New("Error BeginTransaction")).AnyTimes()

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	// Run test
	err = rp.CreateUserAsync("id", "email", "fullname", "phone", "passwordHash")
	assert.Error(t, err)
	assert.EqualError(t, err, "init transactions failed after 5 attempts: Error InitTransactions")

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

	eventsChan := make(chan kev.Event, 1) // buffered supaya gak blocking

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

	eventsChan := make(chan kev.Event, 1)
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

	eventsChan := make(chan kev.Event, 1)
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

	eventsChan := make(chan kev.Event, 1)
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

func TestCreateUserAsync_Success(t *testing.T) {
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

	eventsChan := make(chan kev.Event, 1)
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

func TestCreateUserAsync_Success_EnsureSchemaReturnsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup basic mocks
	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	// Mock ensureSchema success case - GetLatestSchemaRegistery returns nil error (schema exists)
	// This means ensureSchema will return nil without calling CreateAvroSchema
	mockClient.EXPECT().
		GetLatestSchemaMetadata("insert-account--value").
		Return(schrgs.SchemaMetadata{ID: 1, Subject: "insert-account--value", Version: 1}, nil).
		Times(1)

	mockClient.EXPECT().
		GetLatestSchemaMetadata("insert-user--value").
		Return(schrgs.SchemaMetadata{ID: 2, Subject: "insert-user--value", Version: 1}, nil).
		Times(1)

	// Mock Avro serializer
	mockSerializer := mocks.NewMockAvrSerializer(ctrl)
	mockAvr.EXPECT().
		NewGenericSerializer(gomock.Any(), avr.ValueSerde, gomock.Any()).
		Return(mockSerializer, nil)

	// Mock serialization for both account and user
	mockSerializer.EXPECT().
		Serialize("insert-account", gomock.Any()).
		Return([]byte("account-payload"), nil).
		Times(1)

	mockSerializer.EXPECT().
		Serialize("insert-user", gomock.Any()).
		Return([]byte("user-payload"), nil).
		Times(1)

	// Mock Kafka producer
	mockProducer := mocks.NewMockKevProducer(ctrl)

	mockKev.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil)

	eventsChan := make(chan kev.Event, 2)
	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()

	mockProducer.EXPECT().
		InitTransactions(gomock.Any()).
		Return(nil)

	mockProducer.EXPECT().
		BeginTransaction().
		Return(nil)

	// Mock produce messages for both account and user - fix Return() to only have 1 argument
	mockProducer.EXPECT().
		Produce(gomock.Any(), nil).
		Return(nil).
		Times(2)

	// Mock commit transaction
	mockProducer.EXPECT().
		CommitTransaction(gomock.Any()).
		Return(nil)

	// Create repository instance
	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	// Execute the method under test
	err = rp.CreateUserAsync("test-id", "test@email.com", "Test User", "+1234567890", "hashed-password")

	// Assert success
	assert.NoError(t, err)
	close(eventsChan)
}

func TestCreateUserAsync_Success_SchemaNotFoundThenCreatedSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)
	setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

	// Mock: Schema tidak ada (40401 error), kemudian berhasil dibuat
	mockClient.EXPECT().
		GetLatestSchemaMetadata("insert-account--value").
		Return(schrgs.SchemaMetadata{}, errors.New("40401")).
		Times(1)

	mockClient.EXPECT().
		Register("insert-account--value", gomock.Any(), false).
		Return(1, nil). // Berhasil register
		Times(1)

	mockClient.EXPECT().
		GetLatestSchemaMetadata("insert-user--value").
		Return(schrgs.SchemaMetadata{}, errors.New("40401")).
		Times(1)

	mockClient.EXPECT().
		Register("insert-user--value", gomock.Any(), false).
		Return(2, nil). // Berhasil register
		Times(1)

	// Mock Avro serializer
	mockSerializer := mocks.NewMockAvrSerializer(ctrl)
	mockAvr.EXPECT().
		NewGenericSerializer(gomock.Any(), avr.ValueSerde, gomock.Any()).
		Return(mockSerializer, nil)

	mockSerializer.EXPECT().
		Serialize("insert-account", gomock.Any()).
		Return([]byte("account-payload"), nil)

	mockSerializer.EXPECT().
		Serialize("insert-user", gomock.Any()).
		Return([]byte("user-payload"), nil)

	// Mock Kafka producer
	mockProducer := mocks.NewMockKevProducer(ctrl)
	mockKev.EXPECT().
		NewProducer(gomock.Any()).
		Return(mockProducer, nil)

	eventsChan := make(chan kev.Event, 2)
	mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()

	mockProducer.EXPECT().InitTransactions(gomock.Any()).Return(nil)
	mockProducer.EXPECT().BeginTransaction().Return(nil)
	mockProducer.EXPECT().Produce(gomock.Any(), nil).Return(nil).Times(2)
	mockProducer.EXPECT().CommitTransaction(gomock.Any()).Return(nil)

	rp, err := repository.NewPostgresRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
	assert.NoError(t, err)

	err = rp.CreateUserAsync("test-id", "test@email.com", "Test User", "+1234567890", "hashed-password")
	assert.NoError(t, err)

	close(eventsChan)
}
