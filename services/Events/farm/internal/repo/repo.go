package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/concurrent"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

type FarmRepo interface {
	CloseRepo()
	FarmConsumer() kev.KevConsumer
	FarmAddrConsumer() kev.KevConsumer
	DeserializerFarm(topic string, payload []byte) (f models.Farm, _ error)
	DeserializerFarmAddress(topic string, payload []byte) (f models.FarmAddress, _ error)
	UpsertFarmCache(ctx context.Context, farm models.Farm, ops string) error
	UpsertFarmAddressCache(ctx context.Context, addr models.FarmAddress) error
	DeleteFarmCache(ctx context.Context, farmID string, farmAddrID string, farmerID string) error
}

const (
	ConsumerFarmsType = "ConsumerFarmsType"
	ConsumerAddrType  = "ConsumerAddrType"
)

type farmRepo struct {
	srcClient        schrgs.SchrgsClient
	avrDeserializer  avr.AvrDeserializer
	farmConsumer     kev.KevConsumer
	farmAddrConsumer kev.KevConsumer
	farmCache        redis.RedisClient
	stateDB          *pebble.DB
}

type srcAvr struct {
	srClient     schrgs.SchrgsClient
	deserializer avr.AvrDeserializer
}

type farmConsumers struct {
	consumerType string
	consumer     kev.KevConsumer
}

func initSchemaRegistery(
	ctx context.Context,
	sri schrgs.SchemaRegisteryInstance,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[schrgs.SchrgsClient]

		client, err := sri.NewClient(
			sri.NewConfig(
				os.Getenv("SCHEMAREGISTERYADDR"),
			),
		)
		if err != nil {
			res.Error = err
			concurrent.SendResult(ctx, out, res)
			return
		}

		res.Value = client
		concurrent.SendResult(ctx, out, res)
	}()
	return out
}

func initSrcClientAndAvr(
	ctx context.Context,
	avri avr.AvrSerdeInstance,
	srClient <-chan any,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[srcAvr]

		client := <-srClient

		schRes, ok := client.(concurrent.Result[schrgs.SchrgsClient])
		if !ok {
			res.Error = errors.New("wrong type data")
			concurrent.SendResult(ctx, out, res)
			return
		}

		if schRes.Error != nil {
			res.Error = schRes.Error
			concurrent.SendResult(ctx, out, res)
			return
		}

		deseri, err := avri.NewGenericDeserializer(
			schRes.Value,
			avr.ValueSerde,
			avr.NewDeserializerConfig(),
		)
		if err != nil {
			res.Error = err
			concurrent.SendResult(ctx, out, res)
			return
		}

		res.Value.srClient = schRes.Value
		res.Value.deserializer = deseri
		concurrent.SendResult(ctx, out, res)
	}()
	return out
}

func redisClientConn(ctx context.Context, rdi redis.RedisInstance) (redis.RedisClient, error) {
	count := 0
	rdb := redis.NewRedisDB(rdi)
	var errConn error

	for range 5 {

		count++
		rdc, err := rdb.InitRedisClient(
			ctx,
			&redis.FailoverOptions{
				MasterName: os.Getenv("FARM_REDIS_MASTER_NAME"),
				SentinelAddrs: []string{
					os.Getenv("SENTINEL_FARM_REDIS_ADDR"),
					os.Getenv("SENTINEL_FARM_REDIS_ADDR_2"),
				},
				Username: os.Getenv("FARM_REDIS_MASTER_USER_NAME"),
				Password: os.Getenv("FARM_REDIS_MASTER_PASSWORD"),
				DB:       0,
			},
		)

		if err == nil {
			return rdc, nil
		}

		errConn = err
		time.Sleep(time.Second * 2)
	}

	return nil, fmt.Errorf("connection failed after %d attempt: %w", count, errConn)
}

func initFarmCache(ctx context.Context, rdi redis.RedisInstance) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[redis.RedisClient]

		rdc, err := redisClientConn(ctx, rdi)
		if err != nil {
			res.Error = err
			concurrent.SendResult(ctx, out, res)
			return
		}

		res.Value = rdc
		concurrent.SendResult(ctx, out, res)
	}()
	return out
}

func initKafkaFarmConsumer(ctx context.Context, kv kev.Kafka) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[farmConsumers]

		pool := kev.NewKafkaConsumerPool(kv)

		cfgFarm := map[kev.ConfigKeyKafka]string{
			kev.BOOTSTRAP_SERVERS:             os.Getenv("KAKFKABROKER"),
			kev.GROUP_ID:                      "farm-farms-sync-event",
			kev.AUTO_OFFSET_RESET:             "earliest",
			kev.ENABLE_AUTO_COMMIT:            "false",
			kev.PARTITION_ASSIGNMENT_STRATEGY: "cooperative-sticky",
		}

		consumer, err := pool.Consumer(cfgFarm)
		if err != nil {
			res.Error = err
			concurrent.SendResult(ctx, out, res)
			return
		}
		res.Value.consumerType = ConsumerFarmsType
		res.Value.consumer = consumer
		concurrent.SendResult(ctx, out, res)
	}()
	return out
}

func initKafkaFarmAddrConsumer(ctx context.Context, kv kev.Kafka) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[farmConsumers]

		pool := kev.NewKafkaConsumerPool(kv)

		cfgFarm := map[kev.ConfigKeyKafka]string{
			kev.BOOTSTRAP_SERVERS:             os.Getenv("KAKFKABROKER"),
			kev.GROUP_ID:                      "farm-address-sync-event",
			kev.AUTO_OFFSET_RESET:             "earliest",
			kev.ENABLE_AUTO_COMMIT:            "false",
			kev.PARTITION_ASSIGNMENT_STRATEGY: "cooperative-sticky",
		}

		consumer, err := pool.Consumer(cfgFarm)
		if err != nil {
			res.Error = err
			concurrent.SendResult(ctx, out, res)
			return
		}

		res.Value.consumerType = ConsumerAddrType
		res.Value.consumer = consumer
		concurrent.SendResult(ctx, out, res)
	}()
	return out
}

func NewFarmRepo(
	ctx context.Context,
	sri schrgs.SchemaRegisteryInstance,
	avri avr.AvrSerdeInstance,
	kv kev.Kafka,
	rdi redis.RedisInstance,
	stateDB *pebble.DB,
) (fr farmRepo, _ error) {
	fr.stateDB = stateDB

	srCh := initSchemaRegistery(ctx, sri)
	chs := []<-chan any{
		initSrcClientAndAvr(ctx, avri, srCh),
		initKafkaFarmConsumer(ctx, kv),
		initKafkaFarmAddrConsumer(ctx, kv),
		initFarmCache(ctx, rdi),
	}

	for v := range concurrent.FanIn(ctx, chs...) {
		switch res := v.(type) {
		case concurrent.Result[srcAvr]:
			if res.Error != nil {
				return fr, res.Error
			}
			fr.srcClient = res.Value.srClient
			fr.avrDeserializer = res.Value.deserializer
		case concurrent.Result[redis.RedisClient]:
			if res.Error != nil {
				return fr, res.Error
			}

			fr.farmCache = res.Value
		case concurrent.Result[farmConsumers]:
			if res.Error != nil {
				return fr, res.Error
			}
			switch res.Value.consumerType {
			case ConsumerFarmsType:
				fr.farmConsumer = res.Value.consumer
			case ConsumerAddrType:
				fr.farmAddrConsumer = res.Value.consumer
			}
		}
	}

	return fr, nil
}

func (fr farmRepo) CloseRepo() {
	fr.srcClient.Close()
	fr.avrDeserializer.Close()
	fr.farmConsumer.Close()
	fr.farmCache.Close()
}

func (fr farmRepo) FarmConsumer() kev.KevConsumer {
	return fr.farmConsumer
}

func (fr farmRepo) FarmAddrConsumer() kev.KevConsumer {
	return fr.farmAddrConsumer
}

func (fr farmRepo) DeserializerFarm(topic string, payload []byte) (f models.Farm, _ error) {
	if err := fr.avrDeserializer.DeserializeInto(topic, payload, &f); err != nil {
		log.Println(err)
		return f, err
	}

	return f, nil
}

func (fr farmRepo) DeserializerFarmAddress(topic string, payload []byte) (f models.FarmAddress, _ error) {
	if err := fr.avrDeserializer.DeserializeInto(topic, payload, &f); err != nil {
		log.Println(err)
		return f, err
	}

	return f, nil
}

func (fr farmRepo) UpsertFarmCache(
	ctx context.Context,
	farm models.Farm,
	ops string,
) error {
	key := fmt.Sprintf("farm:%s:%s", farm.ID, farm.FarmerID)

	pipe := fr.farmCache.TxPipeline()

	hset := pipe.HSet(ctx, key, farm)
	if hset.Err() != nil {
		return hset.Err()
	}

	expire := pipe.Expire(ctx, key, time.Hour*24)
	if expire.Err() != nil {
		return expire.Err()
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	if ops != "u" {
		go func() {
			_, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			if err := fr.stateDB.Set([]byte(farm.AddressID), []byte(farm.ID), pebble.Sync); err != nil {
				log.Println(err)
				return
			}
		}()
	}

	return nil
}

func (fr farmRepo) UpsertFarmAddressCache(ctx context.Context, addr models.FarmAddress) error {
	var key string

	dataFarmID, closer, err := fr.stateDB.Get([]byte(addr.ID))
	switch {
	case errors.Is(err, pebble.ErrNotFound):
		key = fmt.Sprintf("farm_address:%s:*", addr.ID)
	case err != nil:
		return err
	default:
		defer closer.Close()
	}

	if len(key) <= 0 {
		key = fmt.Sprintf("farm_address:%s:%s", addr.ID, string(dataFarmID))
	}

	if err := fr.stateDB.Delete([]byte(addr.ID), pebble.Sync); err != nil {
		return err
	}

	pipe := fr.farmCache.TxPipeline()

	hset := pipe.HSet(ctx, key, addr)
	if hset.Err() != nil {
		return hset.Err()
	}

	expire := pipe.Expire(ctx, key, time.Hour*24)
	if expire.Err() != nil {
		return expire.Err()
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (fr farmRepo) DeleteFarmCache(ctx context.Context, farmID string, farmAddrID string, farmerID string) error {
	keyFarm := fmt.Sprintf("farm:%s:%s", farmID, farmerID)
	keyFarmAddr := fmt.Sprintf("farm_address:%s:%s", farmAddrID, farmID)

	pipe := fr.farmCache.TxPipeline()

	delFarm := pipe.Del(ctx, keyFarm)
	if delFarm.Err() != nil {
		return delFarm.Err()
	}

	delFarmAddr := pipe.Del(ctx, keyFarmAddr)
	if delFarmAddr.Err() != nil {
		return delFarm.Err()
	}

	return nil
}
