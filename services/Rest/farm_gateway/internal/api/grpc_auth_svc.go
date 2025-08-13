package api

import (
	"context"
	"time"

	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

type GrpcAuthService interface {
	AuthUserRegister(req *pbgen.RegisterUserRequest) (*pbgen.RegisterUserResponse, error)
	AuthUserSignIn(req *pbgen.AuthenticateUserRequest) (*pbgen.AuthenticateUserResponse, error)
	AuthTokenValidate(req *pbgen.TokenValidateRequest) (*pbgen.TokenValidateResponse, error)
}

type grpcService struct {
	authSvc pbgen.AuthServiceClient
}

func NewGrpcService(conn *GrpcClientConn) GrpcAuthService {
	authSvc := pbgen.NewAuthServiceClient(conn)
	return grpcService{authSvc: authSvc}
}

func (s grpcService) AuthUserRegister(req *pbgen.RegisterUserRequest) (*pbgen.RegisterUserResponse, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	res, err := s.authSvc.RegisterUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s grpcService) AuthUserSignIn(req *pbgen.AuthenticateUserRequest) (*pbgen.AuthenticateUserResponse, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	res, err := s.authSvc.AuthenticateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s grpcService) AuthTokenValidate(req *pbgen.TokenValidateRequest) (*pbgen.TokenValidateResponse, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	res, err := s.authSvc.TokenValidate(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
