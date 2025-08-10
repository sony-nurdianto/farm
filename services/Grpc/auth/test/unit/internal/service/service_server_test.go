package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/service"
	"github.com/sony-nurdianto/farm/auth/internal/usecase"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/stretchr/testify/assert"
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

	request := &pbgen.RegisterRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "SomePassword",
	}

	_, err := svc.Register(
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

	request := &pbgen.RegisterRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "SomePassword",
	}

	_, err := svc.Register(
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

	request := &pbgen.RegisterRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "SomePassword",
	}

	_, err := svc.Register(
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

	request := &pbgen.RegisterRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "SomePassword",
	}

	_, err := svc.Register(
		context.Background(),
		request,
	)
	assert.Error(t, err)
}

func TestServiceRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockServiceUsecase(ctrl)
	mockUsecase.EXPECT().
		UserRegister(gomock.Any()).
		Return(
			&pbgen.RegisterResponse{
				Status: "Success",
				Msg:    "Success Register User",
			},
			nil,
		)

	svc := service.NewAuthServiceServer(mockUsecase)

	request := &pbgen.RegisterRequest{
		FullName:    "Sony",
		Email:       "Sony@gmail.com",
		PhoneNumber: "+62851588206",
		Password:    "SomePassword",
	}

	_, err := svc.Register(
		context.Background(),
		request,
	)
	assert.NoError(t, err)
}
