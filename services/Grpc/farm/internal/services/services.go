package services

import (
	"context"
	"io"
	"log"

	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/usescase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FarmServiceServer struct {
	pbgen.UnimplementedFarmServiceServer
	farmUc usescase.FarmUsecase
}

func (fss FarmServiceServer) CreateFarm(stream pbgen.FarmService_CreateFarmServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			log.Println("CreateFarm Services Done")
			return nil
		default:
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Println("createFarm Services Response is done")
				return nil
			}
			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}

			createFarm, err := fss.farmUc.InsertUsers(ctx, msg)
			if err == nil {
				if err := stream.Send(createFarm); err != nil {
					return status.Error(codes.Internal, err.Error())
				}
			}

			if err := stream.Send(createFarm); err != nil {
				return status.Error(codes.Internal, err.Error())
			}

		}
	}
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
