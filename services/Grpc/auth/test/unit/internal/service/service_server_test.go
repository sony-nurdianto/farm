package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/service"
	"github.com/sony-nurdianto/farm/auth/internal/usecase"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServiceRegisterUserIsExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockServiceUsecase(ctrl)
	mockUsecase.EXPECT().
		UserRegister(gomock.Any()).
		Return(
			nil,
			usecase.ErrorUserIsExist,
		)

	svc := service.NewAuthServiceServer(mockUsecase)

	request := &pbgen.RegisterUserRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "SomePassword",
	}

	_, err := svc.RegisterUser(
		context.Background(),
		request,
	)
	assert.Error(t, err)
}

func TestServiceRegisterFailedHashPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockServiceUsecase(ctrl)
	mockUsecase.EXPECT().
		UserRegister(gomock.Any()).
		Return(
			nil,
			usecase.ErrorFailedToHasshPassword,
		)

	svc := service.NewAuthServiceServer(mockUsecase)

	request := &pbgen.RegisterUserRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "SomePassword",
	}

	_, err := svc.RegisterUser(
		context.Background(),
		request,
	)
	assert.Error(t, err)
}

func TestServiceRegisterCreateUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockServiceUsecase(ctrl)
	mockUsecase.EXPECT().
		UserRegister(gomock.Any()).
		Return(
			nil,
			usecase.ErrorRegisterUser,
		)

	svc := service.NewAuthServiceServer(mockUsecase)

	request := &pbgen.RegisterUserRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "SomePassword",
	}

	_, err := svc.RegisterUser(
		context.Background(),
		request,
	)
	assert.Error(t, err)
}

func TestServiceRegisterUnknownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockServiceUsecase(ctrl)
	mockUsecase.EXPECT().
		UserRegister(gomock.Any()).
		Return(
			nil,
			errors.New("Something Wrong"),
		)

	svc := service.NewAuthServiceServer(mockUsecase)

	request := &pbgen.RegisterUserRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "SomePassword",
	}

	_, err := svc.RegisterUser(
		context.Background(),
		request,
	)
	assert.Error(t, err)
}

func TestServiceRegisterSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockServiceUsecase(ctrl)
	mockUsecase.EXPECT().
		UserRegister(gomock.Any()).
		Return(
			&pbgen.RegisterUserResponse{
				Status: "Success",
				Msg:    "Success Register User",
			},
			nil,
		)

	svc := service.NewAuthServiceServer(mockUsecase)

	request := &pbgen.RegisterUserRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "SomePassword",
	}

	_, err := svc.RegisterUser(
		context.Background(),
		request,
	)
	assert.NoError(t, err)
}

func TestServiceUserLoginErrorUserIsNotExsist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockServiceUsecase(ctrl)
	mockUsecase.EXPECT().
		UserSignIn(gomock.Any()).
		Return(nil, usecase.ErrorUserIsNotExsist)

	svc := service.NewAuthServiceServer(mockUsecase)

	request := &pbgen.AuthenticateUserRequest{
		Email:    "Sony@gmail.com",
		Password: "SomePassword",
	}

	_, err := svc.AuthenticateUser(
		context.Background(),
		request,
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, usecase.ErrorUserIsNotExsist.Error())
}

func TestServiceUserLoginErrorPasswordInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockServiceUsecase(ctrl)
	mockUsecase.EXPECT().
		UserSignIn(gomock.Any()).
		Return(nil, usecase.ErrorPasswordIsInvalid)

	svc := service.NewAuthServiceServer(mockUsecase)

	request := &pbgen.AuthenticateUserRequest{
		Email:    "Sony@gmail.com",
		Password: "SomePassword",
	}

	_, err := svc.AuthenticateUser(
		context.Background(),
		request,
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, usecase.ErrorPasswordIsInvalid.Error())
}

func TestServiceUserLoginErrorInternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockServiceUsecase(ctrl)
	mockUsecase.EXPECT().
		UserSignIn(gomock.Any()).
		Return(nil, errors.New("Db Is Invalid"))

	svc := service.NewAuthServiceServer(mockUsecase)

	request := &pbgen.AuthenticateUserRequest{
		Email:    "Sony@gmail.com",
		Password: "SomePassword",
	}

	_, err := svc.AuthenticateUser(
		context.Background(),
		request,
	)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "Db Is Invalid")
}

func TestServiceUserLoginSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockServiceUsecase(ctrl)
	mockUsecase.EXPECT().
		UserSignIn(gomock.Any()).
		Return(
			&pbgen.AuthenticateUserResponse{
				Token:     "Token",
				Status:    "Success",
				Msg:       "Success Register User",
				IssuedAt:  timestamppb.Now(),
				ExpiresAt: timestamppb.New(time.Now().Add(1 * time.Hour)),
			},
			nil,
		)

	svc := service.NewAuthServiceServer(mockUsecase)

	request := &pbgen.AuthenticateUserRequest{
		Email:    "Sony@gmail.com",
		Password: "SomePassword",
	}

	_, err := svc.AuthenticateUser(
		context.Background(),
		request,
	)
	assert.NoError(t, err)
}
