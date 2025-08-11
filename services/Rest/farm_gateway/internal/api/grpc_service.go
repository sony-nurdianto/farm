package api

import (
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

type GrpcService struct {
	authSvc pbgen.AuthServiceClient
}

func NewGrpcService(conn *GrpcClientConn) GrpcService {
	authSvc := pbgen.NewAuthServiceClient(conn)
	return GrpcService{authSvc: authSvc}
}

func (s GrpcService) Services() pbgen.AuthServiceClient {
	return s.authSvc
}
