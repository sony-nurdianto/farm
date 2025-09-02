package repo

import (
	"context"
	"log"
	"time"

	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/models"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/pbgen"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
)

func (fr farmRepo) GetTotalFarms(
	ctx context.Context,
	req *pbgen.GetFarmListRequest,
) (int, error) {
	var total int
	row := fr.farmDB.totalFarmsData.QueryRowContext(ctx, req.GetFarmerId(), req.SearchName)

	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func farmWithAddressScanner(rws pkg.Rows) ([]models.FarmWithAddress, error) {
	var res []models.FarmWithAddress

	defer rws.Close()

	for rws.Next() {
		var farm models.FarmWithAddress

		if err := rws.Scan(
			&farm.Farm.ID,
			&farm.FarmerID,
			&farm.FarmName,
			&farm.FarmType,
			&farm.FarmSize,
			&farm.FarmStatus,
			&farm.Description,
			&farm.Farm.CreatedAt,
			&farm.Farm.UpdatedAt,
			&farm.FarmAddress.ID,
			&farm.Street,
			&farm.Village,
			&farm.SubDistrict,
			&farm.City,
			&farm.Province,
			&farm.PostalCode,
		); err != nil {
			return res, err
		}

		res = append(res, farm)
	}

	return res, nil
}

func (fr farmRepo) GetFarms(
	ctx context.Context,
	req *pbgen.GetFarmListRequest,
) (res []models.FarmWithAddress, _ error) {
	if req.SortOrder == pbgen.SortOrder_SortOrder_ASC {
		rows, err := fr.farmDB.getFarmsAscStmt.QueryContext(ctx, req.GetFarmerId(), req.SearchName, req.Limit, req.Offset)
		if err != nil {
			return res, err
		}

		res, err = farmWithAddressScanner(rows)
		if err != nil {
			return res, err
		}

	}

	if req.SortOrder == pbgen.SortOrder_SortOrder_DESC {
		rows, err := fr.farmDB.getFarmsDescStmt.QueryContext(ctx, req.GetFarmerId(), req.SearchName, req.Limit, req.Offset)
		if err != nil {
			return res, err
		}

		res, err = farmWithAddressScanner(rows)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}

func (fr farmRepo) GetFarmByID(
	ctx context.Context,
	id string,
) (res models.FarmWithAddress, _ error) {
	cache, err := fr.getFarmCache(ctx, id)
	if err == nil {
		log.Println("Return From Cache")
		return cache, nil
	}

	row := fr.farmDB.getFarmByIDStmt.QueryRowContext(ctx, id)

	if err := row.Scan(
		&res.Farm.ID,
		&res.FarmerID,
		&res.FarmName,
		&res.FarmType,
		&res.FarmSize,
		&res.FarmStatus,
		&res.Description,
		&res.Farm.CreatedAt,
		&res.Farm.UpdatedAt,
		&res.FarmAddress.ID,
		&res.Street,
		&res.Village,
		&res.SubDistrict,
		&res.City,
		&res.Province,
		&res.PostalCode,
	); err != nil {
		return res, err
	}

	res.AddressesID = res.FarmAddress.ID

	go func(f models.Farm, a models.FarmAddress) {
		cacheCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := fr.insertFarmCache(cacheCtx, f, a); err != nil {
			log.Println(err)
		}
	}(res.Farm, res.FarmAddress)

	log.Println("Return From Database")

	return res, nil
}
