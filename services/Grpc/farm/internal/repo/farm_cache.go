package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/concurent"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
)

// import (
//
//	"context"
//	"errors"
//	"fmt"
//	"time"
//
//	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/models"
//
// )
//
// var (
//
//	ErrFarmCacheNotExist    = errors.New("farm data not found in cache")
//	ErrFarmAddressNotExsist = errors.New("farm address data not found in cache")
//
// )

func farmCacheToFarmWithAddress(
	ctx context.Context,
	data map[string]string,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurent.Result[models.FarmWithAddress]
		var redisFarm models.RedisFarm

		searchFactory := redis.NewMapStringStringCmd(ctx)
		searchFactory.SetVal(data)

		if err := searchFactory.Scan(&redisFarm); err != nil {
			res.Error = err
			concurent.SendResult(ctx, out, res)
			return
		}

		res.Value = redisFarm.RedisFarmToFarmWithAddress()
		concurent.SendResult(ctx, out, res)
	}()
	return out
}

func (fr farmRepo) getFarmCache(
	ctx context.Context,
	search string,
	farmerID string,
	limit int,
	offset int,
	asc bool,
) (res []models.FarmWithAddress, err error) {
	var query string
	factoryFarmerID := strings.ReplaceAll(farmerID, "-", "\\-")
	if search == "" || search == "*" {
		query = fmt.Sprintf("@farmer_id:{%s}", factoryFarmerID)
	} else {
		query = fmt.Sprintf("@farm_name:*%s* @farmer_id:{%s}", search, factoryFarmerID)
	}
	farmCache := fr.farmCache.FTSearchWithArgs(
		ctx,
		"farm_idx",
		query,
		&redis.FTSearchOptions{
			SortBy: []redis.FTSearchSortBy{
				{FieldName: "created_at", Asc: asc},
			},
			Limit:       limit,
			LimitOffset: offset,
		},
	)

	if farmCache.Err() != nil {
		return res, farmCache.Err()
	}

	if len(farmCache.Val().Docs) <= 0 {
		return res, nil
	}

	chs := make([]<-chan any, len(farmCache.Val().Docs))

	for i, v := range farmCache.Val().Docs {
		chs[i] = farmCacheToFarmWithAddress(ctx, v.Fields)
	}

	for v := range concurent.FanIn(ctx, chs...) {
		resCh, ok := v.(concurent.Result[models.FarmWithAddress])
		if !ok {
			return res, errors.New("wrong data type")
		}

		if resCh.Error != nil {
			return res, resCh.Error
		}

		fmt.Println(resCh.Value)

		res = append(res, resCh.Value)
	}

	return res, nil
}

func sendDataFarmCache(
	ctx context.Context,
	serializer avr.AvrSerializer,
	producer kev.KevProducer,
	farm models.AvrFarm,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurent.Result[struct{}]

		topic := "insert-farm-cache"

		farmPayload, err := serializer.Serialize(topic, farm)
		if err != nil {
			res.Error = err
			concurent.SendResult(ctx, out, res)
			return
		}

		msg := kev.MessageKafka{
			TopicPartition: kev.KafkaTopicPartition{
				Topic:     &topic,
				Partition: kev.KafkaPartitionAny,
			},
			Value: farmPayload,
		}.Factory()

		if err := producer.Produce(&msg, nil); err != nil {
			res.Error = err
			concurent.SendResult(ctx, out, res)
			return
		}

		res.Value = struct{}{}

		concurent.SendResult(ctx, out, res)
	}()
	return out
}

func (fr farmRepo) insertFarmCache(farms ...models.FarmWithAddress) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	chs := make([]<-chan any, len(farms))

	for i, farm := range farms {
		avrFarm := models.AvrFarm{
			ID:               farm.Farm.ID,
			FarmerID:         farm.FarmerID,
			FarmName:         farm.FarmName,
			FarmType:         farm.FarmType,
			FarmSize:         farm.FarmSize,
			FarmStatus:       farm.FarmStatus,
			Description:      &farm.Description,
			CreatedAt:        farm.Farm.CreatedAt,
			UpdatedAt:        farm.Farm.UpdatedAt,
			AddressID:        farm.AddressesID,
			Street:           farm.Street,
			Village:          farm.Village,
			SubDistrict:      farm.SubDistrict,
			City:             farm.City,
			Province:         farm.Province,
			PostalCode:       farm.PostalCode,
			AddressCreatedAt: farm.FarmAddress.CreatedAt,
			AddressUpdatedAt: farm.FarmAddress.UpdatedAt,
		}

		chs[i] = sendDataFarmCache(ctx, fr.avrSerializer, fr.farmProducer, avrFarm)
	}

	for v := range concurent.FanIn(ctx, chs...) {
		res, ok := v.(concurent.Result[struct{}])
		if !ok {
			return errors.New("wrong data type")
		}
		if res.Error != nil {
			log.Println(res.Error)
			return res.Error
		}
	}

	return nil
}
