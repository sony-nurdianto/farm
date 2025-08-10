package unit_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/auth/internal/entity"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/usecase"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUseCaseUserRegisterUserExsist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{Email: "test@gmail.com"}, nil)

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn)

	req := &pbgen.RegisterRequest{
		FullName:    "test",
		Email:       "test@gmail.com",
		PhoneNumber: "+47545687898",
		Password:    "Something",
	}

	_, err := uc.UserRegister(req)
	assert.Error(t, err)
	assert.EqualError(t, err, usecase.ErrorUserIsExist.Error())
}

func TestUserRegister_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, sql.ErrConnDone)

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn)

	req := &pbgen.RegisterRequest{
		FullName:    "test",
		Email:       "test@gmail.com",
		PhoneNumber: "+47545687898",
		Password:    "Something",
	}

	_, err := uc.UserRegister(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "sql: connection is already closed")
}

func TestUserRegister_HashPasswordError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, sql.ErrNoRows)

	mocksPassEn.EXPECT().
		HashPassword(gomock.Any()).
		Return("", errors.New("Failed To HashPassword"))

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn)

	req := &pbgen.RegisterRequest{
		FullName:    "test",
		Email:       "test@gmail.com",
		PhoneNumber: "+47545687898",
		Password:    "Something",
	}

	_, err := uc.UserRegister(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed To HashPassword")
}

func TestUserRegister_CreateUserAsyncdError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, sql.ErrNoRows)

	mocksPassEn.EXPECT().
		HashPassword(gomock.Any()).
		Return("HashPassword", nil)

	mockAuthRepo.EXPECT().
		CreateUserAsync(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(errors.New("Failed Create User"))

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn)

	req := &pbgen.RegisterRequest{
		FullName:    "test",
		Email:       "test@gmail.com",
		PhoneNumber: "+47545687898",
		Password:    "Something",
	}

	_, err := uc.UserRegister(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "Failed To CreateUserAsync: Failed Create User")
}

func TestUserRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, sql.ErrNoRows)

	mocksPassEn.EXPECT().
		HashPassword(gomock.Any()).
		Return("HashPassword", nil)

	mockAuthRepo.EXPECT().
		CreateUserAsync(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil)

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn)

	req := &pbgen.RegisterRequest{
		FullName:    "test",
		Email:       "test@gmail.com",
		PhoneNumber: "+47545687898",
		Password:    "Something",
	}

	out, err := uc.UserRegister(req)
	assert.NoError(t, err)
	assert.Equal(t, out.Msg, "Success Create User")
	assert.Equal(t, out.Status, "Success")
}
