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

func updateFarmSendMsg(
	ctx context.Context,
	stream grpc.BidiStreamingClient[pbgen.UpdateFarmsRequest, pbgen.UpdateFarmsResponse],
	dataRequest []models.UpdateFarmWithAddr,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)

		var res concurrent.Result[struct{}]
		for _, req := range dataRequest {

			var farm *pbgen.UpdateFarmData
			var address *pbgen.UpdateFarmAddressData

			if req.Farm != nil {
				farm = &pbgen.UpdateFarmData{
					Id:          req.Farm.ID,
					FarmName:    req.Farm.FarmName,
					FarmType:    req.Farm.FarmType,
					FarmStatus:  req.Farm.FarmStatus,
					FarmSize:    req.Farm.FarmSize,
					Description: req.Farm.Description,
				}
			}

			if req.Address != nil {
				address = &pbgen.UpdateFarmAddressData{
					Id:          req.Address.ID,
					Street:      req.Address.Street,
					Village:     req.Address.Village,
					SubDistrict: req.Address.SubDistrict,
					City:        req.Address.City,
					Province:    req.Address.Province,
					PostalCode:  req.Address.PostalCode,
				}
			}

			msg := &pbgen.UpdateFarmsRequest{
				Farm:    farm,
				Address: address,
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

func updateFarmRecvMsg(
	ctx context.Context,
	stream grpc.BidiStreamingClient[pbgen.UpdateFarmsRequest, pbgen.UpdateFarmsResponse],
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var recv []*pbgen.UpdateFarmsResponse
		var res concurrent.Result[[]*pbgen.UpdateFarmsResponse]
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

func (s grpcFarmService) UpdateFarmOrAddress(
	ctx context.Context,
	data []models.UpdateFarmWithAddr,
) ([]*pbgen.UpdateFarmsResponse, error) {
	var results []*pbgen.UpdateFarmsResponse
	stream, err := s.farmSvc.UpdateFarms(ctx)
	if err != nil {
		return nil, err
	}

	chs := []<-chan any{
		updateFarmSendMsg(ctx, stream, data),
		updateFarmRecvMsg(ctx, stream),
	}

	for v := range concurrent.FanIn(ctx, chs...) {
		switch res := v.(type) {
		case concurrent.Result[struct{}]:
			if res.Error != nil {
				return nil, res.Error
			}
		case concurrent.Result[[]*pbgen.UpdateFarmsResponse]:
			if res.Error != nil {
				return nil, res.Error
			}

			results = res.Value
		}
	}

	return results, nil
}
