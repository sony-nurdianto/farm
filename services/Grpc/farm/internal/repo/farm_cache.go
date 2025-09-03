package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/models"
)

var (
	ErrFarmCacheNotExist    = errors.New("farm data not found in cache")
	ErrFarmAddressNotExsist = errors.New("farm address data not found in cache")
)

func (fr farmRepo) getFarmCache(ctx context.Context, key string, farmerID string) (res models.FarmWithAddress, _ error) {
	var farm models.Farm
	var addr models.FarmAddress

	if err := fr.farmCache.HGetAll(ctx, fmt.Sprintf("farm:%s:%s", key, farmerID)).Scan(&farm); err != nil {
		return res, err
	}

	if farm == (models.Farm{}) {
		return res, ErrFarmCacheNotExist
	}

	if err := fr.farmCache.HGetAll(ctx, fmt.Sprintf("farm_address:%s:%s", farm.AddressesID, farm.ID)).Scan(&addr); err != nil {
		return res, nil
	}

	if addr == (models.FarmAddress{}) {
		return res, ErrFarmAddressNotExsist
	}

	res.Farm = farm
	res.FarmAddress = addr

	return res, nil
}

func (fr farmRepo) insertFarmCache(ctx context.Context, farm models.Farm, addr models.FarmAddress) error {
	pipe := fr.farmCache.TxPipeline()

	insertFarmKey := fmt.Sprintf("farm:%s:%s", farm.ID, farm.FarmerID)
	insertFarm := pipe.HSet(ctx, insertFarmKey, farm)
	if insertFarm.Err() != nil {
		return insertFarm.Err()
	}

	insertAddrKey := fmt.Sprintf("farm_address:%s:%s", addr.ID, farm.ID)
	insertAddr := pipe.HSet(ctx, insertAddrKey, addr)
	if insertAddr.Err() != nil {
		return insertAddr.Err()
	}

	if err := pipe.Expire(ctx, insertFarmKey, time.Hour*24).Err(); err != nil {
		return err
	}

	if err := pipe.Expire(ctx, insertAddrKey, time.Hour*24).Err(); err != nil {
		return err
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	return nil
}
