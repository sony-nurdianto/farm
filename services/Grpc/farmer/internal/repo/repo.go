package repo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/concurent"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/models"

	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

type FarmerRepo interface {
	GetUsersByIDFromCache(ctx context.Context, id string) (farmer models.Users, _ error)
	UpdateUserAsync(ctx context.Context, users *models.UpdateUsers) error
}

type farmerRepo struct {
	sriClient      schrgs.SchrgsClient
	avrSerializer  avr.AvrSerializer
	farmerProudcer kev.KevProducer
	farmerCache    redis.RedisClient
}

func send(
	ctx context.Context,
	send chan any,
	recv any,
) {
	select {
	case <-ctx.Done():
		return
	case send <- recv:
	}
}

func initSchemaRegistery(ctx context.Context, sri schrgs.SchemaRegisteryInstance) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurent.Result[schrgs.SchrgsClient]
		src, err := sri.NewClient(
			sri.NewConfig(os.Getenv("SCHEMAREGISTERYADDR")),
		)
		if err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}

		res.Value = src
		send(ctx, out, res)
	}()

	return out
}

type schemaRegistryPair struct {
	Serializer avr.AvrSerializer
	Client     schrgs.SchrgsClient
}

func schemaAndSerializerPipe(
	ctx context.Context,
	avri avr.AvrSerdeInstance,
	clientChan <-chan any,
) <-chan any {
	out := make(chan any, 1)

	go func() {
		defer close(out)
		var res concurent.Result[schemaRegistryPair]

		client := <-clientChan

		schRes, ok := client.(concurent.Result[schrgs.SchrgsClient])
		if !ok {
			res.Error = errors.New("wrong type data")
			send(ctx, out, res)
			return
		}

		if schRes.Error != nil {
			res.Error = schRes.Error
			send(ctx, out, res)
			return
		}

		seri, err := avri.NewGenericSerializer(
			schRes.Value,
			avr.ValueSerde,
			avr.NewSerializerConfig(),
		)
		if err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}
		res.Value.Serializer = seri
		res.Value.Client = schRes.Value
		send(ctx, out, res)
	}()
	return out
}

func initFarmerProducer(ctx context.Context, kv kev.Kafka) <-chan any {
	out := make(chan any, 1)

	go func() {
		defer close(out)
		var res concurent.Result[kev.KevProducer]

		pool := kev.NewKafkaProducerPool(kv)
		cfg := map[kev.ConfigKeyKafka]string{
			kev.BOOTSTRAP_SERVERS:  os.Getenv("KAKFKABROKER"),
			kev.ACKS:               "all",
			kev.ENABLE_IDEMPOTENCE: "true",
			kev.COMPRESION_TYPE:    "snappy",
			kev.RETRIES:            "5",
			kev.RETRY_BACKOFF_MS:   "100",
			kev.LINGER_MS:          "5",
			kev.MAX_IN_FLIGHT_REQUESTS_PER_CONNECTION: "5",
			kev.TRANSACTIONAL_ID:                      "",
		}

		producer, err := pool.Producer(cfg)
		if err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}

		res.Value = producer
		send(ctx, out, res)
	}()

	return out
}

func initRedisDatabae(ctx context.Context, rdi redis.RedisInstance) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurent.Result[redis.RedisClient]

		rdb := redis.NewRedisDB(rdi)
		rdc, err := rdb.InitRedisClient(context.Background(), &redis.FailoverOptions{
			MasterName:    os.Getenv("FARMER_REDIS_MASTER_NAME"),
			SentinelAddrs: []string{os.Getenv("SENTINEL_FARMER_REDIS_ADDR")},
			Username:      os.Getenv("FARMER_REDIS_MASTER_USER_NAME"),
			Password:      os.Getenv("FARMER_REDIS_MASTER_PASSWORD"),
			DB:            0,
		})
		if err != nil {
			res.Error = err
			send(ctx, out, res)
			return
		}

		res.Value = rdc
		send(ctx, out, res)
	}()

	return out
}

func NewFarmerRepo(
	ctx context.Context,
	avri avr.AvrSerdeInstance,
	kv kev.Kafka,
	sri schrgs.SchemaRegisteryInstance,
	rdi redis.RedisInstance,
) (fr farmerRepo, err error) {
	opsCtx, done := context.WithTimeout(ctx, time.Second*30)
	defer done()

	src := initSchemaRegistery(opsCtx, sri)
	chs := []<-chan any{
		schemaAndSerializerPipe(opsCtx, avri, src),
		initFarmerProducer(opsCtx, kv),
		initRedisDatabae(opsCtx, rdi),
	}

	for v := range concurent.FanIn(opsCtx, chs...) {
		switch res := v.(type) {
		case concurent.Result[schemaRegistryPair]:
			if res.Error != nil {
				return fr, res.Error
			}
			fr.sriClient = res.Value.Client
			fr.avrSerializer = res.Value.Serializer
		case concurent.Result[kev.KevProducer]:
			if res.Error != nil {
				return fr, res.Error
			}
			fr.farmerProudcer = res.Value
		case concurent.Result[redis.RedisClient]:
			if res.Error != nil {
				return fr, res.Error
			}
			fr.farmerCache = res.Value
		}
	}

	return fr, nil
}

func (fr farmerRepo) GetUsersByIDFromCache(ctx context.Context, id string) (user models.Users, _ error) {
	hkey := fmt.Sprintf("users:%s", id)

	err := fr.farmerCache.HGetAll(ctx, hkey).Scan(&user)
	if err != nil {
		return user, nil
	}
	if user == (models.Users{}) {
		return user, errors.New("user is not existed")
	}
	return user, nil

}
