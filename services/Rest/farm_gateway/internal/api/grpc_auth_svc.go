package api

import (
	"context"
	"time"

	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

//go:generate mockgen -source=../pbgen/auth_grpc.pb.go -destination=../../test/mocks/mock_auth_grpc.pb.go -package=mocks
//go:generate mockgen -source=grpc_auth_svc.go -destination=../../test/mocks/mock_grpc_auth_svc.go -package=mocks

type GrpcAuthService interface {
	AuthUserRegister(req *pbgen.RegisterUserRequest) (*pbgen.RegisterUserResponse, error)
	AuthUserSignIn(req *pbgen.AuthenticateUserRequest) (*pbgen.AuthenticateUserResponse, error)
	AuthTokenValidate(req *pbgen.TokenValidateRequest) (*pbgen.TokenValidateResponse, error)
}

type grpcAuthService struct {
	authSvc pbgen.AuthServiceClient
}

func NewGrpcService(svc pbgen.AuthServiceClient) GrpcAuthService {
	return grpcAuthService{authSvc: svc}
}

func (s grpcAuthService) AuthUserRegister(req *pbgen.RegisterUserRequest) (*pbgen.RegisterUserResponse, error) {
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

func (s grpcAuthService) AuthUserSignIn(req *pbgen.AuthenticateUserRequest) (*pbgen.AuthenticateUserResponse, error) {
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

func (s grpcAuthService) AuthTokenValidate(req *pbgen.TokenValidateRequest) (*pbgen.TokenValidateResponse, error) {
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
