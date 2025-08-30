package repo

import (
	"context"
	"time"

	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
)

func changeFarmAddreses(
	ctx context.Context, tx pkg.Stmt, address *models.UpdateFarmAddress,
) (res models.FarmAddress, _ error) {
	row := tx.QueryRowContext(
		ctx,
		address.Street,
		address.Village,
		address.SubDistrict,
		address.City,
		address.Province,
		address.PostalCode,
		time.Now().UTC(), // updated_at
		address.ID,
	)

	if err := row.Err(); err != nil {
		return res, err
	}

	if err := row.Scan(
		&res.ID,
		&res.Street,
		&res.Village,
		&res.SubDistrict,
		&res.City,
		&res.Province,
		&res.PostalCode,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		return res, err
	}

	return res, nil
}

func changeFarm(
	ctx context.Context, tx pkg.Stmt, farm *models.UpdateFarm,
) (res models.Farm, _ error) {
	row := tx.QueryRowContext(
		ctx,
		farm.FarmName,
		farm.FarmType,
		farm.FarmSize,
		farm.FarmStatus,
		farm.Description,
		time.Now().UTC(),
	)

	if err := row.Err(); err != nil {
		return res, err
	}

	if err := row.Scan(
		&res.ID,
		&res.FarmerID,
		&res.FarmName,
		&res.FarmType,
		&res.FarmSize,
		&res.FarmStatus,
		&res.Description,
		&res.AddressesID,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		return res, err
	}

	return res, nil
}

func (fr farmRepo) UpdateFarm(
	ctx context.Context, opts *pkg.TxOpts, farm *models.UpdateFarm, address *models.UpdateFarmAddress,
) (*models.Farm, *models.FarmAddress, error) {
	tx, err := fr.farmDB.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	defer tx.Rollback()

	var farmRes *models.Farm
	var farmAddrRes *models.FarmAddress

	if farm != nil {
		txFarmStmt := tx.Stmt(fr.farmDB.updateFarmStmt)
		farm, err := changeFarm(ctx, txFarmStmt, farm)
		if err != nil {
			return nil, nil, err
		}

		farmRes = &farm
	}

	if address != nil {
		txAddrStmt := tx.Stmt(fr.farmDB.updateFarmAddresStmt)
		addRes, err := changeFarmAddreses(ctx, txAddrStmt, address)
		if err != nil {
			return nil, nil, err
		}
		farmAddrRes = &addRes
	}

	tx.Commit()

	return farmRes, farmAddrRes, nil
}
