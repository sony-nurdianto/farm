package repo

import (
	"context"

	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/models"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/pbgen"
)

func (fr farmRepo) GetTotalFarms(
	ctx context.Context,
	req *pbgen.GetFarmListRequest,
) (int, error) {
	var total int
	row := fr.farmDB.totalFarmsData.QueryRowContext(ctx, req.GetFarmerId(), req.SearchName, req.Limit, req.Offset)

	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
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

		defer rows.Close()

		for rows.Next() {
			var farm models.FarmWithAddress

			if err := rows.Scan(
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
	}

	if req.SortOrder == pbgen.SortOrder_SortOrder_DESC {
		rows, err := fr.farmDB.getFarmsDescStmt.QueryContext(ctx, req.GetFarmerId(), req.SearchName, req.Limit, req.Offset)
		if err != nil {
			return res, err
		}

		defer rows.Close()

		for rows.Next() {
			var farm models.FarmWithAddress

			if err := rows.Scan(
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
	}
	return res, nil
}
