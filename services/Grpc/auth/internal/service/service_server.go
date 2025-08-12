package service

import (
	"context"
	"errors"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceServer struct {
	pbgen.UnimplementedAuthServiceServer
	serviceUsecase usecase.ServiceUsecase
}

func NewAuthServiceServer(uc usecase.ServiceUsecase) *AuthServiceServer {
	return &AuthServiceServer{serviceUsecase: uc}
}

func handleRegisterError(err error) error {
	switch {
	case err == usecase.ErrorUserIsExist:
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, usecase.ErrorFailedToHasshPassword):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, usecase.ErrorRegisterUser):
		return status.Error(codes.Internal, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (ass *AuthServiceServer) RegisterUser(
	ctx context.Context,
	in *pbgen.RegisterUserRequest,
) (*pbgen.RegisterUserResponse, error) {
	res, err := ass.serviceUsecase.UserRegister(in)
	if err != nil {
		return nil, handleRegisterError(err)
	}

	return res, nil
}

func (ass *AuthServiceServer) AuthenticateUser(
	ctx context.Context,
	in *pbgen.AuthenticateUserRequest,
) (*pbgen.AuthenticateUserResponse, error) {
	return nil, nil
}
