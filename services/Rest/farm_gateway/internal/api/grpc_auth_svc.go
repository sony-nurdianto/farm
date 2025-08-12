package api

import (
	"context"
	"time"

	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

type GrpcAuthService interface {
	AuthUserRegister(req *pbgen.RegisterUserRequest) (*pbgen.RegisterUserResponse, error)
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
