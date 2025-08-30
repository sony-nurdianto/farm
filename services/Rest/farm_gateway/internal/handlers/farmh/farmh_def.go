package farmh

import "github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/api"

type farmHandler struct {
	grpcFarmSvc api.GrpcFarmService
}

func NewFarmHandler(grpcSvc api.GrpcFarmService) farmHandler {
	return farmHandler{
		grpcFarmSvc: grpcSvc,
	}
}
