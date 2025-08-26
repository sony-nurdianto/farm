package services

import (
	"context"

	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/pbgen"
)

type FarmServiceServer struct {
	pbgen.UnimplementedFarmServiceServer
}

func (fss FarmServiceServer) CreateFarm(stream pbgen.FarmService_CreateFarmServer) error {
	return nil
}

func (fss FarmServiceServer) GetFarmByID(ctx context.Context, in *pbgen.GetFarmByIDRequest) (*pbgen.GetFarmByIDResponse, error) {
	return nil, nil
}

func (fss FarmServiceServer) GetFarmList(in *pbgen.GetFarmListRequest, stream pbgen.FarmService_GetFarmListServer) error {
	return nil
}

func (fss FarmServiceServer) UpdateFarms(stream pbgen.FarmService_UpdateFarmsServer) error {
	return nil
}

func (fss FarmServiceServer) DeleteFarm(ctx context.Context, in *pbgen.DeleteFarmRequest) (*pbgen.DeleteFarmResponse, error) {
	return nil, nil
}
