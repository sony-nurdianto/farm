package repo

import (
	"context"

	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
)

func (fr farmRepo) insertFarmAddress(
	ctx context.Context, tx pkg.Stmt, address models.FarmAddress,
) (res models.FarmAddress, _ error) {
	row := tx.QueryRowContext(
		ctx,
		address.ID,
		address.Street,
		address.Village,
		address.SubDistrict,
		address.City,
		address.Province,
		address.PostalCode,
		address.CreatedAt,
		address.UpdatedAt,
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

func (fr farmRepo) insertFarm(
	ctx context.Context, tx pkg.Stmt, farm models.Farm,
) (res models.Farm, _ error) {
	row := tx.QueryRowContext(
		ctx,
		farm.ID,
		farm.FarmerID,
		farm.FarmName,
		farm.FarmType,
		farm.FarmSize,
		farm.FarmStatus,
		farm.Description,
		farm.AddressesID,
		farm.CreatedAt,
		farm.UpdatedAt,
	)
	if row.Err() != nil {
		return res, row.Err()
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

func (fr farmRepo) CreateFarm(
	ctx context.Context,
	opts pkg.TxOpts,
	farm models.Farm,
	farmAddr models.FarmAddress,
) (res models.FarmWithAddress, _ error) {
	tx, err := fr.farmDB.db.BeginTx(ctx, &opts)
	if err != nil {
		return res, err
	}

	defer tx.Rollback()

	txAddrStmt := tx.Stmt(fr.farmDB.createFarmAddressStmt)
	addrRes, err := fr.insertFarmAddress(ctx, txAddrStmt, farmAddr)
	if err != nil {
		return res, err
	}

	txFarmStmt := tx.Stmt(fr.farmDB.createFarmStmt)
	farmRes, err := fr.insertFarm(ctx, txFarmStmt, farm)
	if err != nil {
		return res, err
	}

	if err := tx.Commit(); err != nil {
		return res, err
	}

	res.Farm = farmRes
	res.FarmAddress = addrRes
	return res, nil
}
