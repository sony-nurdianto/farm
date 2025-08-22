package service

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/pbgen"
	"go.opentelemetry.io/otel/metric"
)

type farmerServiceServer struct {
	pbgen.UnimplementedFarmerServiceServer
	tracer trace.Tracer
	meter  metric.Meter
}

func NewFarmerServiceServer(trc trace.Tracer, mtr metric.Meter) *farmerServiceServer {
	return &farmerServiceServer{
		tracer: trc,
		meter:  mtr,
	}
}

func (fss farmerServiceServer) FarmerProfile(
	ctx context.Context, in *pbgen.FarmerProfileRequest,
) (*pbgen.FarmerProfileResponse, error) {
	return nil, nil
}

func (fss farmerServiceServer) UpdateFarmerProfile(
	ctx context.Context, in *pbgen.UpdateFarmerProfileRequest,
) (*pbgen.UpdateFarmerProfileResponse, error) {
	return nil, nil
}
