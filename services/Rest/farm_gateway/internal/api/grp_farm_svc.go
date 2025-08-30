package api

import (
	"context"
	"io"
	"log"

	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/concurrent"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
	"google.golang.org/grpc"
)

type GrpcFarmService interface {
	CreateFarm(ctx context.Context, dataRequest []models.CreateFarm) ([]*pbgen.CreateFarmResponse, error)
}

type grpcFarmService struct {
	farmSvc pbgen.FarmServiceClient
}

func NewGrpcFarmService(svc pbgen.FarmServiceClient) GrpcFarmService {
	return grpcFarmService{
		farmSvc: svc,
	}
}

func createFarmSendMsg(
	ctx context.Context,
	stream grpc.BidiStreamingClient[pbgen.CreateFarmRequest, pbgen.CreateFarmResponse],
	dataRequest []models.CreateFarm,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[struct{}]
		for _, req := range dataRequest {
			msg := &pbgen.CreateFarmRequest{
				Farm: &pbgen.CreateFarm{
					FarmerId:    req.FarmerID,
					FarmName:    req.FarmName,
					FarmType:    req.FarmType,
					FarmSize:    req.FarmSize,
					FarmStatus:  req.FarmStatus,
					Description: req.Description,
				},
				Address: &pbgen.CreateFarmAddress{
					Street:      req.Address.Street,
					Village:     req.Address.Village,
					SubDistrict: req.Address.SubDistrict,
					City:        req.Address.City,
					Province:    req.Address.Province,
					PostalCode:  req.Address.PostalCode,
				},
			}

			if err := stream.Send(msg); err != nil {
				res.Error = err
				concurrent.SendResult(ctx, out, res)
				return
			}
		}

		if err := stream.CloseSend(); err != nil {
			res.Error = err
			concurrent.SendResult(ctx, out, res)
			return
		}

		res.Value = struct{}{}
		concurrent.SendResult(ctx, out, res)
	}()
	return out
}

func createFarmRecvMsg(
	ctx context.Context,
	stream grpc.BidiStreamingClient[pbgen.CreateFarmRequest, pbgen.CreateFarmResponse],
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var recv []*pbgen.CreateFarmResponse
		var res concurrent.Result[[]*pbgen.CreateFarmResponse]

		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Println("Request Done")
				res.Value = recv
				concurrent.SendResult(ctx, out, res)
				return
			}

			if err != nil {
				res.Error = err
				concurrent.SendResult(ctx, out, res)
				return
			}

			recv = append(recv, msg)
		}
	}()
	return out
}

func (s grpcFarmService) CreateFarm(
	ctx context.Context,
	dataRequest []models.CreateFarm,
) ([]*pbgen.CreateFarmResponse, error) {
	var results []*pbgen.CreateFarmResponse

	stream, err := s.farmSvc.CreateFarm(ctx)
	if err != nil {
		return results, err
	}

	chs := []<-chan any{
		createFarmSendMsg(ctx, stream, dataRequest),
		createFarmRecvMsg(ctx, stream),
	}

	for v := range concurrent.FanIn(ctx, chs...) {
		switch res := v.(type) {
		case concurrent.Result[struct{}]:
			if res.Error != nil {
				return nil, res.Error
			}
		case concurrent.Result[[]*pbgen.CreateFarmResponse]:
			if res.Error != nil {
				return nil, res.Error
			}
			results = res.Value
		}
	}

	return results, nil
}
