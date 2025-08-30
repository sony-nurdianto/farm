package api

import (
	"context"

	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

type GrpcFarmService interface {
	CreateFarm(ctx context.Context, dataRequest []models.CreateFarm) ([]*pbgen.CreateFarmResponse, error)
	UpdateFarmOrAddress(ctx context.Context, data []models.UpdateFarmWithAddr) ([]*pbgen.UpdateFarmsResponse, error)
}

type grpcFarmService struct {
	farmSvc pbgen.FarmServiceClient
}

func NewGrpcFarmService(svc pbgen.FarmServiceClient) GrpcFarmService {
	return grpcFarmService{
		farmSvc: svc,
	}
}
