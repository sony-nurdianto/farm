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

func NewFarmServiceServer(uc usescase.FarmUsecase) FarmServiceServer {
	return FarmServiceServer{
		farmUc: uc,
	}
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

			createFarm := fss.farmUc.InsertUsers(ctx, msg)

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
	ctx := stream.Context()

	totalFarm, err := fss.farmUc.GetTotalFarms(ctx, in)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	err = stream.Send(&pbgen.GetFarmListResponse{
		Total: int32(totalFarm),
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	farms, err := fss.farmUc.GetFarms(ctx, in)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, v := range farms {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := stream.Send(v); err != nil {
				return status.Error(codes.Internal, err.Error())
			}

		}
	}

	return nil
}

func (fss FarmServiceServer) UpdateFarms(stream pbgen.FarmService_UpdateFarmsServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Println("updateFarm Or UpdateAddress services is Done")
				return nil
			}
			if err != nil {
				return status.Error(codes.Internal, err.Error())
			}

			if msg.Farm == nil && msg.Address == nil {
				return status.Error(codes.InvalidArgument, "at least farm or farm address have value")
			}

			updateFarm := fss.farmUc.UpdateUsers(ctx, msg)
			if err := stream.Send(updateFarm); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}
	}
}

func (fss FarmServiceServer) DeleteFarm(ctx context.Context, in *pbgen.DeleteFarmRequest) (*pbgen.DeleteFarmResponse, error) {
	return nil, nil
}
