package repo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sony-nurdianto/farm/services/Events/farm/internal/concurrent"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

func InitSchemaRegistery(
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

func redisClientConn(
	ctx context.Context,
	rdi redis.RedisInstance,
	unstableResp3 bool,
	protocol int,
) (redis.RedisClient, error) {
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
				Username:      os.Getenv("FARM_REDIS_MASTER_USER_NAME"),
				Password:      os.Getenv("FARM_REDIS_MASTER_PASSWORD"),
				DB:            0,
				Protocol:      protocol,
				UnstableResp3: unstableResp3,
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

func InitFarmCache(
	ctx context.Context,
	unstableResp3 bool,
	protocol int,
	rdi redis.RedisInstance,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[redis.RedisClient]

		rdc, err := redisClientConn(ctx, rdi, unstableResp3, protocol)
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
