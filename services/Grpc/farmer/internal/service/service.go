package service

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/models"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/usecase"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
	"go.opentelemetry.io/otel/metric"
)

type farmerServiceServer struct {
	pbgen.UnimplementedFarmerServiceServer
	tracer        trace.Tracer
	meter         metric.Meter
	farmerUsecase usecase.FarmerUsecase
}

func NewFarmerServiceServer(trc trace.Tracer, mtr metric.Meter, farmerUsecase usecase.FarmerUsecase) *farmerServiceServer {
	return &farmerServiceServer{
		tracer:        trc,
		meter:         mtr,
		farmerUsecase: farmerUsecase,
	}
}

func (fss farmerServiceServer) FarmerProfile(
	ctx context.Context, in *pbgen.FarmerProfileRequest,
) (*pbgen.FarmerProfileResponse, error) {
	farmer, err := fss.farmerUsecase.GetUserByID(ctx, in.GetId())
	if err == nil {
		return farmer, nil
	}

	if err == redis.RedisNil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("user with id %s is not exist", in.GetId()))
	}

	return nil, err
}

func (fss farmerServiceServer) UpdateFarmerProfile(
	ctx context.Context, in *pbgen.UpdateFarmerProfileRequest,
) (*pbgen.UpdateFarmerProfileResponse, error) {
	userUpdate := &models.UpdateUsers{
		ID:       in.GetId(),
		FullName: in.FullName,
		Email:    in.Email,
		Phone:    in.Phone,
	}

	if _, err := fss.farmerUsecase.UpdateUser(ctx, userUpdate); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &pbgen.UpdateFarmerProfileResponse{
		Status: "Success",
		Msg:    "Sucess Update Farmer Profile",
	}

	return res, nil
}
