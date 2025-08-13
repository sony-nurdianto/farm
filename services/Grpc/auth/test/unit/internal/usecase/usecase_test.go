package unit_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/o1egl/paseto"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/token"
	"github.com/sony-nurdianto/farm/auth/internal/entity"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/usecase"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUseCaseUserRegisterUserExsist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{Email: "test@gmail.com"}, nil)

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.RegisterUserRequest{
		FullName:    "test",
		Email:       "test@gmail.com",
		PhoneNumber: "+47545687898",
		Password:    "Something",
	}

	_, err := uc.UserRegister(req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, usecase.ErrorUserIsExist)
}

func TestUserRegister_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, sql.ErrConnDone)

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.RegisterUserRequest{
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
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, sql.ErrNoRows)

	mocksPassEn.EXPECT().
		HashPassword(gomock.Any()).
		Return("", errors.New("Failed To HashPassword"))

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.RegisterUserRequest{
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
	mocksTokhan := mocks.NewMockTokhan(ctrl)

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

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.RegisterUserRequest{
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
	mocksTokhan := mocks.NewMockTokhan(ctrl)

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

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.RegisterUserRequest{
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

func TestUserSignIn_ErrorUserNotExsist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, sql.ErrNoRows)

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.AuthenticateUserRequest{
		Email:    "test@gmail.com",
		Password: "Something",
	}

	_, err := uc.UserSignIn(req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, usecase.ErrorUserIsNotExsist)
}

func TestUserSignIn_ErrorGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, errors.New("Db is Not Defined"))

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.AuthenticateUserRequest{
		Email:    "test@gmail.com",
		Password: "Something",
	}

	_, err := uc.UserSignIn(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "Db is Not Defined")
}

func TestUserSignIn_ErrVerifyPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, nil)

	mocksPassEn.EXPECT().
		VerifyPassword(gomock.Any(), gomock.Any()).
		Return(false, errors.New("Error VerifyPassword"))

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.AuthenticateUserRequest{
		Email:    "test@gmail.com",
		Password: "Something",
	}

	_, err := uc.UserSignIn(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "Error VerifyPassword")
}

func TestUserSignIn_ErrPasswordInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, nil)

	mocksPassEn.EXPECT().
		VerifyPassword(gomock.Any(), gomock.Any()).
		Return(false, nil)

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.AuthenticateUserRequest{
		Email:    "test@gmail.com",
		Password: "Something",
	}

	_, err := uc.UserSignIn(req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, usecase.ErrorPasswordIsInvalid)
}

func TestUserSignIn_ErrCreateWebToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, nil)

	mocksPassEn.EXPECT().
		VerifyPassword(gomock.Any(), gomock.Any()).
		Return(true, nil)

	mocksTokhan.EXPECT().
		CreateWebToken(gomock.Any()).
		Return("", errors.New("Unexpected Error when create web token"))

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.AuthenticateUserRequest{
		Email:    "test@gmail.com",
		Password: "Something",
	}

	_, err := uc.UserSignIn(req)
	assert.Error(t, err)
	assert.EqualError(t, err, "Unexpected Error when create web token")
}

func TestUserSignIn_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	mockAuthRepo.EXPECT().
		GetUserByEmail(gomock.Any()).
		Return(entity.Users{}, nil)

	mocksPassEn.EXPECT().
		VerifyPassword(gomock.Any(), gomock.Any()).
		Return(true, nil)

	mocksTokhan.EXPECT().
		CreateWebToken(gomock.Any()).
		Return("Token", nil)

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.AuthenticateUserRequest{
		Email:    "test@gmail.com",
		Password: "Something",
	}

	out, err := uc.UserSignIn(req)
	assert.NoError(t, err)
	assert.Equal(t, out.Token, "Token")
	assert.Equal(t, out.Msg, "User Authenticated Success Login. Welcome !")
	assert.Equal(t, out.Status, "Success")
}

func TestTokenValidateExperied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	mocksTokhan.EXPECT().
		VerifyWebToken(gomock.Any()).
		Return(paseto.JSONToken{
			Issuer:     "auth",
			Expiration: time.Now().Add(time.Hour * -1),
			Subject:    "user-id",
		}, token.ErrTokenExperied)

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.TokenValidateRequest{
		Token: "token",
	}

	res, err := uc.TokenValidate(req)
	assert.NoError(t, err)
	assert.Nil(t, res.Isuer)
	assert.Nil(t, res.Subject)
	assert.Nil(t, res.ExpiresAt)
	assert.False(t, res.Valid)
	assert.Equal(t, res.Msg, "Token Experied")
}

func TestTokenValidateErrorVerifyWebToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	mocksTokhan.EXPECT().
		VerifyWebToken(gomock.Any()).
		Return(paseto.JSONToken{}, token.ErrDecryptFailed)

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.TokenValidateRequest{
		Token: "token",
	}

	res, err := uc.TokenValidate(req)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.EqualError(t, err, token.ErrDecryptFailed.Error())
}

func TestTokenValidateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepo(ctrl)
	mocksPassEn := mocks.NewMockPassEncrypt(ctrl)
	mocksTokhan := mocks.NewMockTokhan(ctrl)

	expRes := paseto.JSONToken{
		Issuer:     "auth",
		Expiration: time.Now().Add(time.Hour * 1),
		Subject:    "user-id",
	}

	mocksTokhan.EXPECT().
		VerifyWebToken(gomock.Any()).
		Return(expRes, nil)

	uc := usecase.NewServiceUsecase(mockAuthRepo, mocksPassEn, mocksTokhan)

	req := &pbgen.TokenValidateRequest{
		Token: "token",
	}

	res, err := uc.TokenValidate(req)
	assert.NoError(t, err)
	assert.Equal(t, *res.Isuer, expRes.Issuer)
	assert.Equal(t, *res.Subject, expRes.Subject)
	assert.Equal(t, res.ExpiresAt, timestamppb.New(expRes.Expiration))
	assert.True(t, res.Valid)
	assert.Equal(t, res.Msg, "Token Is Valid")
}
