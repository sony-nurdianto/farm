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

func TestGetUser(t *testing.T) {
	t.Run("GetUserByEmail Error Scan", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev, mockProducer := setupBasicMocks(ctrl)

		setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockKev, mockProducer)

		eventsChan := make(chan kev.Event, 1)
		mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()

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

		close(eventsChan)
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockAvr, mockKev, mockProducer := setupBasicMocks(ctrl)

		setupRepositoryMocks(mockPgI, mockPgDb, mockStmt, mockSchrgs, mockClient, mockKev, mockProducer)

		eventsChan := make(chan kev.Event, 1)
		mockProducer.EXPECT().Events().Return(eventsChan).AnyTimes()

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

		close(eventsChan)
	})
}
