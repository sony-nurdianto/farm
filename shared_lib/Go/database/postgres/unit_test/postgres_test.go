package pkg

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres"
	"github.com/stretchr/testify/mock"
)

type MockPostgres struct {
	mock.Mock
}

func (mp *MockPostgres) Open(driverName, dataSourceName string) (pg postgres.PostgresDatabase, _ error) {
	args := mp.Called(driverName, dataSourceName)

	db, _ := args.Get(0).(postgres.PostgresDatabase)
	return db, args.Error(1)
}

type MockDatabase struct {
	mock.Mock
}

func (md *MockDatabase) SetMaxOpenConns(n int) {
}

func (md *MockDatabase) SetMaxIdleConns(n int) {
}

func (md *MockDatabase) SetConnMaxLifetime(d time.Duration) {
}

func (md *MockDatabase) PingContext(ctx context.Context) error {
	args := md.Called(ctx)
	return args.Error(0)
}

func (md *MockDatabase) Close() error {
	args := md.Called()
	return args.Error(0)
}

func TestOpenPostgres_Error(t *testing.T) {
	mockInstance := new(MockPostgres)

	mockInstance.On("Open", "postgres", "something").Return(nil, errors.New("Failed To Open Postgres"))

	_, err := postgres.OpenPostgres("something", mockInstance)
	if err == nil {
		t.Error("Expected Error when failed to Open Postgres")
	}
}

func TestOpenPostgres_PingError(t *testing.T) {
	mockInstance := new(MockPostgres)
	mockDb := new(MockDatabase)

	mockDb.On("SetMaxOpenConns", mock.Anything).Return()
	mockDb.On("SetMaxIdleConns", mock.Anything).Return()
	mockDb.On("SetConnMaxLifetime", mock.Anything).Return()
	mockDb.On("PingContext", mock.Anything).Return(errors.New("ping failed no response"))
	mockInstance.On("Open", "postgres", "something").Return(mockDb, nil)

	_, err := postgres.OpenPostgres("something", mockInstance)
	if err == nil {
		t.Error("Expected Error PingContext but got nil")
	}
}
