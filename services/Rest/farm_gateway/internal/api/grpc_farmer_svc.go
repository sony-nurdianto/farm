package api

import (
	"context"

	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

type GrpcFarmerService interface {
	FarmerProfile(ctx context.Context, req *pbgen.FarmerProfileRequest) (*pbgen.FarmerProfileResponse, error)
	ProfileFarmerUpdate(ctx context.Context, req *pbgen.UpdateFarmerProfileRequest) (*pbgen.UpdateFarmerProfileResponse, error)
}

type grpcFarmerService struct {
	farmerSvc pbgen.FarmerServiceClient
}

func NewGrpcFarmerService(svc pbgen.FarmerServiceClient) GrpcFarmerService {
	return grpcFarmerService{farmerSvc: svc}
}

func (s grpcFarmerService) FarmerProfile(ctx context.Context, req *pbgen.FarmerProfileRequest) (*pbgen.FarmerProfileResponse, error) {
	res, err := s.farmerSvc.FarmerProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s grpcFarmerService) ProfileFarmerUpdate(ctx context.Context, req *pbgen.UpdateFarmerProfileRequest) (*pbgen.UpdateFarmerProfileResponse, error) {
	res, err := s.farmerSvc.UpdateFarmerProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
