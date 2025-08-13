package authh

import "github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/api"

type authHandler struct {
	grpcAuthSvc api.GrpcAuthService
}

func NewAuthHandler(grpcSvc api.GrpcAuthService) authHandler {
	return authHandler{
		grpcAuthSvc: grpcSvc,
	}
}
