package usescase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/models"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/repo"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
)

type FarmUsecase interface {
	InsertUsers(ctx context.Context, req *pbgen.CreateFarmRequest) (*pbgen.CreateFarmResponse, error)
}

type farmUsecase struct {
	repo repo.FarmRepo
}

func NewFarmUsecase(r repo.FarmRepo) farmUsecase {
	return farmUsecase{
		repo: r,
	}
}

func (fu farmUsecase) InsertUsers(ctx context.Context, req *pbgen.CreateFarmRequest) (*pbgen.CreateFarmResponse, error) {
	txOpts := pkg.TxOpts{
		Isolation: pkg.LevelSerializable,
		ReadOnly:  false,
	}

	fAddrID := uuid.NewString()

	farmAddr := models.FarmAddress{
		ID:          fAddrID,
		Street:      req.Address.GetStreet(),
		Village:     req.Address.GetVillage(),
		SubDistrict: req.Address.GetSubDistrict(),
		City:        req.Address.GetCity(),
		Province:    req.Address.GetProvince(),
		PostalCode:  req.Address.GetPostalCode(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	farm := models.Farm{
		ID:          uuid.NewString(),
		FarmerID:    req.Farm.GetFarmerId(),
		FarmName:    req.Farm.GetFarmName(),
		FarmType:    req.Farm.GetFarmType(),
		FarmSize:    req.Farm.GetFarmSize(),
		FarmStatus:  req.Farm.GetFarmStatus(),
		Description: req.Farm.GetDescription(),
		AddressesID: fAddrID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	res := new(pbgen.CreateFarmResponse)

	users, err := fu.repo.CreateFarm(
		ctx, txOpts, farm, farmAddr,
	)
	if err != nil {
		res.Status = "Error"
		res.Msg = err.Error()

		return res, nil
	}

	res.FarmId = users.Farm.ID
	res.FarmName = users.FarmName
	res.AddressId = users.AddressesID
	res.Status = "Success"
	res.Msg = "Success Create Users"

	return res, nil
}
