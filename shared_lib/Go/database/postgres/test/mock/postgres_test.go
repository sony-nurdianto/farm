package mock

import (
	"errors"
	"testing"

	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOpenPostgres_Error(t *testing.T) {
	mockInstance := new(MockInstance)
	mockDB := new(MockDatabase)

	mockInstance.On("Open", "postgres", "uri").Return(mockDB, errors.New("Failed to open Connection"))
	mockDB.On("SetMaxOpenConns", 30).Return()
	mockDB.On("SetMaxIdleConns", 10).Return()
	mockDB.On("SetConnMaxLifetime", mock.AnythingOfType("time.Duration")).Return()
	mockDB.On("PingContext", mock.Anything).Return(nil)

	db, err := pkg.OpenPostgres("uri", mockInstance)
	assert.Error(t, err)
	assert.Nil(t, db)
	mockInstance.AssertExpectations(t)
}
