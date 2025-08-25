package redis

import (
	"context"
	"fmt"
)

type RedisDatabase interface {
	InitRedisClient(ctx context.Context, opt *FailoverOptions) (RedisClient, error)
}

type rdb struct {
	rdi RedisInstance
}

func NewRedisDB(i RedisInstance) rdb {
	return rdb{
		rdi: i,
	}
}

func (d rdb) InitRedisClient(ctx context.Context, opt *FailoverOptions) (RedisClient, error) {
	rdc := d.rdi.NewFailoverClient(opt)

	ping := rdc.Ping(ctx)
	pong, err := ping.Result()
	if err != nil {
		return nil, err
	}
	fmt.Println("Ping Response: ", pong)

	return NewRedisClient(rdc), nil
}
