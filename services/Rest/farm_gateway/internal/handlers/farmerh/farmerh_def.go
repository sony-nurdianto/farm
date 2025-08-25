package farmerh

import "github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/api"

type farmerHandler struct {
	grpcFarmerSvc api.GrpcFarmerService
}

func NewFarmerHandler(grpcSvc api.GrpcFarmerService) farmerHandler {
	return farmerHandler{
		grpcFarmerSvc: grpcSvc,
	}
}
