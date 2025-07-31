package mock

import (
	"context"
	"time"

	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/stretchr/testify/mock"
)

type MockInstance struct {
	mock.Mock
}

func (mi *MockInstance) Open(driverName string, dataSourceName string) (pkg.PostgresDatabase, error) {
	args := mi.Called(driverName, dataSourceName)
	return args.Get(0).(pkg.PostgresDatabase), args.Error(1)
}

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) SetMaxOpenConns(n int) {
	m.Called(n)
}

func (m *MockDatabase) SetMaxIdleConns(n int) {
	m.Called(n)
}

func (m *MockDatabase) SetConnMaxLifetime(d time.Duration) {
	m.Called(d)
}

func (m *MockDatabase) PingContext(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDatabase) Prepare(query string) (pkg.Stmt, error) {
	args := m.Called(query)
	return args.Get(0).(pkg.Stmt), args.Error(1)
}

func (m *MockDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}
