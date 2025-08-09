package unit_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	t.Run("GetUserByEmail Error Scan", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)

		// Setup repository mocks
		setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

		mockRow := mocks.NewMockRow(ctrl)

		rp, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
		assert.NoError(t, err)

		mockStmt.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any()).
			Return(mockRow)

		mockRow.EXPECT().
			Scan(gomock.Any()).
			Return(errors.New("Failed To Scan Data"))

		_, err = rp.GetUserByEmail("test@email.com")
		assert.Error(t, err)
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev := setupBasicMocks(ctrl)

		// Setup repository mocks
		setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient)

		mockRow := mocks.NewMockRow(ctrl)

		rp, err := repository.NewAuthRepo(mockSchrgs, mockPgI, mockAvr, mockKev)
		assert.NoError(t, err)

		mockStmt.EXPECT().
			QueryRowContext(gomock.Any(), gomock.Any()).
			Return(mockRow)

		mockRow.EXPECT().
			Scan(gomock.Any()).
			Return(nil)

		_, err = rp.GetUserByEmail("test@email.com")
		assert.NoError(t, err)
	})
}
