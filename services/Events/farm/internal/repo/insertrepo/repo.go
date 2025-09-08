package insertrepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sony-nurdianto/farm/services/Events/farm/internal/concurrent"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/models"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/repo"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

type InsertFarmCacheRepo interface {
	CloseRepo()
	Consumer() kev.KevConsumer
	CreateFarmIndex(ctx context.Context) error
	InsertFarmCache(ctx context.Context, farm models.InsertFarm) error
	FarmDeserializer(payload []byte) (fa models.InsertFarm, _ error)
}

type insertFarmCacheRepo struct {
	srcClient       schrgs.SchrgsClient
	avrDeserializer avr.AvrDeserializer
	farmConsumer    kev.KevConsumer
	farmCache       redis.RedisClient
}

type avrSrc struct {
	srClient     schrgs.SchrgsClient
	deserializer avr.AvrDeserializer
}

var ErrIndexAlredyExsist = errors.New("index is alredy exsist")

func initAvrAndSrcClient(
	ctx context.Context,
	avri avr.AvrSerdeInstance,
	srcClient <-chan any,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[avrSrc]

		client, ok := <-srcClient
		if !ok {
			log.Println("Chanel Is Close")
			return
		}
		schRes, ok := client.(concurrent.Result[schrgs.SchrgsClient])
		if !ok {
			res.Error = errors.New("wrong data type")
			concurrent.SendResult(ctx, out, res)
			return
		}
		if schRes.Error != nil {
			res.Error = schRes.Error
			concurrent.SendResult(ctx, out, res)
			return
		}

		deseri, err := avri.NewGenericDeserializer(
			schRes.Value, avr.ValueSerde, avr.NewDeserializerConfig(),
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

func initKafkaInsertFarmConsumer(ctx context.Context, kv kev.Kafka) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[kev.KevConsumer]

		pool := kev.NewKafkaConsumerPool(kv)

		cfg := map[kev.ConfigKeyKafka]string{
			kev.BOOTSTRAP_SERVERS:             os.Getenv("KAKFKABROKER"),
			kev.GROUP_ID:                      "farm-insert-farms-cache-event",
			kev.AUTO_OFFSET_RESET:             "earliest",
			kev.ENABLE_AUTO_COMMIT:            "false",
			kev.PARTITION_ASSIGNMENT_STRATEGY: "cooperative-sticky",
		}

		consumer, err := pool.Consumer(cfg)
		if err != nil {
			res.Error = err
			concurrent.SendResult(ctx, out, res)
			return
		}

		res.Value = consumer
		concurrent.SendResult(ctx, out, res)
	}()
	return out
}

func NewInsertFarmCacheRepo(
	ctx context.Context,
	sri schrgs.SchemaRegisteryInstance,
	avri avr.AvrSerdeInstance,
	rdi redis.RedisInstance,
	kv kev.Kafka,
) (r insertFarmCacheRepo, _ error) {
	srCh := repo.InitSchemaRegistery(ctx, sri)
	chs := []<-chan any{
		initAvrAndSrcClient(ctx, avri, srCh),
		initKafkaInsertFarmConsumer(ctx, kv),
		repo.InitFarmCache(ctx, true, 2, rdi),
	}

	for v := range concurrent.FanIn(ctx, chs...) {
		switch res := v.(type) {
		case concurrent.Result[avrSrc]:
			if res.Error != nil {
				return r, res.Error
			}
			r.srcClient = res.Value.srClient
			r.avrDeserializer = res.Value.deserializer
		case concurrent.Result[kev.KevConsumer]:
			if res.Error != nil {
				return r, res.Error
			}
			r.farmConsumer = res.Value
		case concurrent.Result[redis.RedisClient]:
			if res.Error != nil {
				return r, res.Error
			}
			r.farmCache = res.Value

		}
	}

	return r, nil
}

func (r insertFarmCacheRepo) CloseRepo() {
	r.srcClient.Close()
	r.avrDeserializer.Close()
	r.farmConsumer.Close()
	r.farmCache.Close()
}

func (r insertFarmCacheRepo) CreateFarmIndex(
	ctx context.Context,
) error {
	schemas := []*redis.FieldSchema{
		{FieldName: "id", FieldType: redis.SearchFieldTypeTag},
		{FieldName: "farmer_id", FieldType: redis.SearchFieldTypeTag},
		{FieldName: "address_id", FieldType: redis.SearchFieldTypeTag},
		{FieldName: "farm_name", FieldType: redis.SearchFieldTypeText, Weight: 2.0},
		{FieldName: "farm_type", FieldType: redis.SearchFieldTypeText, Weight: 1.0},
		{FieldName: "farm_status", FieldType: redis.SearchFieldTypeText, Weight: 1.0},
		{FieldName: "description", FieldType: redis.SearchFieldTypeText, Weight: 0.5},
		{FieldName: "farm_size", FieldType: redis.SearchFieldTypeNumeric},
		{FieldName: "postal_code", FieldType: redis.SearchFieldTypeText},
		{FieldName: "street", FieldType: redis.SearchFieldTypeText},
		{FieldName: "village", FieldType: redis.SearchFieldTypeText},
		{FieldName: "sub_district", FieldType: redis.SearchFieldTypeText},
		{FieldName: "city", FieldType: redis.SearchFieldTypeText},
		{FieldName: "province", FieldType: redis.SearchFieldTypeText},
		{FieldName: "photo_url", FieldType: redis.SearchFieldTypeText},
		{FieldName: "created_at", FieldType: redis.SearchFieldTypeNumeric},
		{FieldName: "updated_at", FieldType: redis.SearchFieldTypeNumeric},
		{FieldName: "address_created_at", FieldType: redis.SearchFieldTypeNumeric},
		{FieldName: "address_updated_at", FieldType: redis.SearchFieldTypeNumeric},
	}

	createIdx := r.farmCache.FTCreate(ctx, "farm_idx", &redis.FTCreateOptions{
		OnHash: true,
		Prefix: []any{"farm:"},
	}, schemas...)

	if createIdx.Err() != nil {
		if strings.Contains(createIdx.Err().Error(), "Index already exists") {
			return ErrIndexAlredyExsist
		}

		return createIdx.Err()
	}

	return nil
}

func (r insertFarmCacheRepo) Consumer() kev.KevConsumer {
	return r.farmConsumer
}

func (r insertFarmCacheRepo) FarmDeserializer(payload []byte) (fa models.InsertFarm, _ error) {
	if err := r.avrDeserializer.DeserializeInto("insert-farm-cache", payload, &fa); err != nil {
		log.Println(err)
		return fa, err
	}

	log.Println(fa)

	return fa, nil
}

func (r insertFarmCacheRepo) InsertFarmCache(
	ctx context.Context,
	farm models.InsertFarm,
) error {
	pipe := r.farmCache.TxPipeline()
	key := fmt.Sprintf("farm:%s:%s", farm.ID, farm.FarmerID)

	hSet := pipe.HSet(ctx, key, farm)
	if hSet.Err() != nil {
		return hSet.Err()
	}

	experied := pipe.Expire(ctx, key, time.Hour*24)
	if experied.Err() != nil {
		return experied.Err()
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	return nil
}
