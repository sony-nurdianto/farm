package syncrepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/concurrent"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/models"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/repo"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

type SyncFarmRepo interface {
	CloseRepo()
	FarmConsumer() kev.KevConsumer
	FarmAddrConsumer() kev.KevConsumer
	DeserializerFarm(topic string, payload []byte) (f models.AvrFarm, _ error)
	DeserializerFarmAddress(topic string, payload []byte) (f models.AvrFarmAddress, _ error)
	UpsertFarmCache(ctx context.Context, farm models.AvrFarm, ops string) error
	UpsertFarmAddressCache(ctx context.Context, addr models.AvrFarmAddress) error
	DeleteFarmCache(ctx context.Context, farmID string, farmerID string) error
}

const (
	ConsumerFarmsType = "ConsumerFarmsType"
	ConsumerAddrType  = "ConsumerAddrType"
)

type syncFarmRepo struct {
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

func initSrcClientAndAvr(
	ctx context.Context,
	avri avr.AvrSerdeInstance,
	srClient <-chan any,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[srcAvr]

		client, ok := <-srClient
		if !ok {
			log.Println("Chanel is Close")
			return
		}

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

func NewSyncFarmRepo(
	ctx context.Context,
	sri schrgs.SchemaRegisteryInstance,
	avri avr.AvrSerdeInstance,
	kv kev.Kafka,
	rdi redis.RedisInstance,
	stateDB *pebble.DB,
) (fr syncFarmRepo, _ error) {
	fr.stateDB = stateDB

	srCh := repo.InitSchemaRegistery(ctx, sri)
	chs := []<-chan any{
		initSrcClientAndAvr(ctx, avri, srCh),
		initKafkaFarmConsumer(ctx, kv),
		initKafkaFarmAddrConsumer(ctx, kv),
		repo.InitFarmCache(ctx, false, 3, rdi),
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

func (fr syncFarmRepo) CloseRepo() {
	fr.srcClient.Close()
	fr.avrDeserializer.Close()
	fr.farmConsumer.Close()
	fr.farmCache.Close()
}

func (fr syncFarmRepo) FarmConsumer() kev.KevConsumer {
	return fr.farmConsumer
}

func (fr syncFarmRepo) FarmAddrConsumer() kev.KevConsumer {
	return fr.farmAddrConsumer
}

func (fr syncFarmRepo) DeserializerFarm(topic string, payload []byte) (f models.AvrFarm, _ error) {
	if err := fr.avrDeserializer.DeserializeInto(topic, payload, &f); err != nil {
		log.Println(err)
		return f, err
	}

	return f, nil
}

func (fr syncFarmRepo) DeserializerFarmAddress(topic string, payload []byte) (f models.AvrFarmAddress, _ error) {
	if err := fr.avrDeserializer.DeserializeInto(topic, payload, &f); err != nil {
		log.Println(err)
		return f, err
	}

	return f, nil
}

func (fr syncFarmRepo) UpsertFarmCache(
	ctx context.Context,
	farm models.AvrFarm,
	ops string,
) error {
	key := fmt.Sprintf("farm:%s:%s", farm.ID, farm.FarmerID)

	pipe := fr.farmCache.TxPipeline()

	hset := pipe.HSet(ctx, key, farm.ToFarm())
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

	stateValue := fmt.Sprintf("%s:%s", farm.ID, farm.FarmerID)
	if err := fr.stateDB.Set([]byte(farm.AddressID), []byte(stateValue), pebble.Sync); err != nil {
		return err
	}

	return nil
}

func (fr syncFarmRepo) UpsertFarmAddressCache(ctx context.Context, addr models.AvrFarmAddress) error {
	var farmID string
	var farmerID string

	for range 5 {
		ids, closer, err := fr.stateDB.Get([]byte(addr.ID))
		if err == nil {
			data := strings.Split(string(ids), ":")
			farmID = data[0]
			farmerID = data[1]
			closer.Close()
			break
		}

		if errors.Is(err, pebble.ErrNotFound) {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		return err

	}

	key := fmt.Sprintf("farm:%s:%s", farmID, farmerID)
	hset := fr.farmCache.HSet(ctx, key, addr.ToFarmAddress())
	if hset.Err() != nil {
		return hset.Err()
	}

	if err := fr.stateDB.Delete([]byte(addr.ID), pebble.Sync); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (fr syncFarmRepo) DeleteFarmCache(ctx context.Context, farmID string, farmerID string) error {
	keyFarm := fmt.Sprintf("farm:%s:%s", farmID, farmerID)

	hset := fr.farmCache.Del(ctx, keyFarm)
	if hset.Err() != nil {
		return hset.Err()
	}

	return nil
}
