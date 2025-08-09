package unit_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPostgres_OpenPostgres(t *testing.T) {
	t.Run("Failed To OpenPostgres", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgi := mocks.NewMockPostgresInstance(ctrl)

		mockPgi.EXPECT().
			Open(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("Failed to PostgresDatabase"))

		_, err := pkg.OpenPostgres("someaddress", mockPgi)
		assert.Error(t, err)
		assert.EqualError(t, err, "Failed to PostgresDatabase")
	})

	t.Run("Failed PingContext", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgi := mocks.NewMockPostgresInstance(ctrl)
		mockPgDB := mocks.NewMockPostgresDatabase(ctrl)

		mockPgi.EXPECT().
			Open(gomock.Any(), gomock.Any()).
			Return(mockPgDB, nil)

		mockPgDB.EXPECT().
			SetMaxOpenConns(gomock.Any()).
			Return()

		mockPgDB.EXPECT().
			SetMaxIdleConns(gomock.Any()).
			Return()

		mockPgDB.EXPECT().
			SetConnMaxLifetime(gomock.Any()).
			Return()

		mockPgDB.EXPECT().
			PingContext(gomock.Any()).
			Return(errors.New("Failed PingContext"))

		_, err := pkg.OpenPostgres("someaddress", mockPgi)
		assert.Error(t, err)
		assert.EqualError(t, err, "Failed PingContext")
	})

	t.Run("Sucess OpenDb", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgi := mocks.NewMockPostgresInstance(ctrl)
		mockPgDB := mocks.NewMockPostgresDatabase(ctrl)

		mockPgi.EXPECT().
			Open(gomock.Any(), gomock.Any()).
			Return(mockPgDB, nil)

		mockPgDB.EXPECT().
			SetMaxOpenConns(gomock.Any()).
			Return()

		mockPgDB.EXPECT().
			SetMaxIdleConns(gomock.Any()).
			Return()

		mockPgDB.EXPECT().
			SetConnMaxLifetime(gomock.Any()).
			Return()

		mockPgDB.EXPECT().
			PingContext(gomock.Any()).
			Return(nil)

		_, err := pkg.OpenPostgres("someaddress", mockPgi)
		assert.NoError(t, err)
	})

	t.Run("Prepare Test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgi := mocks.NewMockPostgresInstance(ctrl)
		mockPgDB := mocks.NewMockPostgresDatabase(ctrl)

		mockPgi.EXPECT().
			Open(gomock.Any(), gomock.Any()).
			Return(mockPgDB, nil)

		mockPgDB.EXPECT().
			SetMaxOpenConns(gomock.Any()).
			Return()

		mockPgDB.EXPECT().
			SetMaxIdleConns(gomock.Any()).
			Return()

		mockPgDB.EXPECT().
			SetConnMaxLifetime(gomock.Any()).
			Return()

		mockPgDB.EXPECT().
			PingContext(gomock.Any()).
			Return(nil)

		pg, err := pkg.OpenPostgres("someaddress", mockPgi)
		assert.NoError(t, err)

		mockPgDB.EXPECT().
			Prepare(gomock.Any()).
			Return(nil, errors.New("Failed To Prepare Statement"))

		_, err = pg.Prepare("SomeStatement")
		assert.Error(t, err)
	})
}
