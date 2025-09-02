package api

import (
	"context"
	"io"

	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

func (s grpcFarmService) GetFarmByID(
	ctx context.Context,
	farmID string,
) (res models.Farm, _ error) {
	req := &pbgen.GetFarmByIDRequest{
		Id: farmID,
	}

	farm, err := s.farmSvc.GetFarmByID(ctx, req)
	if err != nil {
		return res, err
	}
	res = models.Farm{
		ID:          farm.Farm.Id,
		FarmerID:    farm.Farm.FarmerId,
		FarmName:    farm.Farm.FarmName,
		FarmType:    farm.Farm.FarmType,
		FarmSize:    farm.Farm.FarmSize,
		FarmStatus:  farm.Farm.FarmStatus,
		Description: farm.Farm.Description,
		Addresses: models.FarmAddress{
			ID:          farm.Farm.Address.Id,
			Street:      farm.Farm.Address.Street,
			Village:     farm.Farm.Address.Village,
			SubDistrict: farm.Farm.Address.SubDistrict,
			City:        farm.Farm.Address.City,
			Province:    farm.Farm.Address.Province,
			PostalCode:  farm.Farm.Address.PostalCode,
		},
		CreatedAt: farm.Farm.CreatedAt.AsTime().UTC(),
		UpdatedAt: farm.Farm.UpdatedAt.AsTime().UTC(),
	}

	return res, nil
}

func (s grpcFarmService) GetFarms(
	ctx context.Context,
	farmerID string,
	dataRequest models.GetFarmsRequest,
) (res models.GetFarmsResponse, _ error) {
	req := &pbgen.GetFarmListRequest{
		FarmerId:   farmerID,
		SearchName: dataRequest.SearchName,
		SortOrder:  dataRequest.SortOrder.ProtoSortOrder(),
		Limit:      int32(dataRequest.Limit),
		Offset:     int32(dataRequest.Offset),
	}

	stream, err := s.farmSvc.GetFarmList(ctx, req)
	if err != nil {
		return res, err
	}

	streamCtx := stream.Context()

	for {
		select {
		case <-streamCtx.Done():
			return res, nil
		default:
			msg, err := stream.Recv()
			if err == io.EOF {
				return res, nil
			}

			if err != nil {
				return res, err
			}

			if msg.Total != nil {
				res.Total = int(msg.GetTotal())
			}

			if msg.Farms != nil {
				farm := models.Farm{
					ID:          msg.Farms.Id,
					FarmerID:    msg.Farms.FarmerId,
					FarmName:    msg.Farms.FarmName,
					FarmType:    msg.Farms.FarmType,
					FarmSize:    msg.Farms.FarmSize,
					FarmStatus:  msg.Farms.FarmStatus,
					Description: msg.Farms.Description,
					Addresses: models.FarmAddress{
						ID:          msg.Farms.Address.Id,
						Street:      msg.Farms.Address.Street,
						Village:     msg.Farms.Address.Village,
						SubDistrict: msg.Farms.Address.SubDistrict,
						City:        msg.Farms.Address.City,
						Province:    msg.Farms.Address.Province,
						PostalCode:  msg.Farms.Address.PostalCode,
					},
					CreatedAt: msg.Farms.CreatedAt.AsTime().UTC(),
					UpdatedAt: msg.Farms.UpdatedAt.AsTime().UTC(),
				}

				res.Data = append(res.Data, farm)
			}
		}
	}
}
