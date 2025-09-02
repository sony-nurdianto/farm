package redis

import (
	"context"
	"time"
)

type RedisClient interface {
	HSet(ctx context.Context, key string, values ...any) *IntCmd
	HGet(ctx context.Context, key string, field string) *StringCmd
	HGetAll(ctx context.Context, key string) *MapStringStringCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd
	TxPipeline() Pipeliner
	Del(ctx context.Context, keys ...string) *IntCmd
	Ping(ctx context.Context) StatusCmd
	Close() error
}

type rdc struct {
	client *Client
}

func NewRedisClient(c *Client) *rdc {
	return &rdc{
		client: c,
	}
}

func (c *rdc) Ping(ctx context.Context) StatusCmd {
	return c.client.Ping(ctx)
}

func (c *rdc) HSet(ctx context.Context, key string, values ...any) *IntCmd {
	return c.client.HSet(ctx, key, values...)
}

func (c *rdc) HGet(ctx context.Context, key string, field string) *StringCmd {
	return c.client.HGet(ctx, key, field)
}

func (c *rdc) HGetAll(ctx context.Context, key string) *MapStringStringCmd {
	return c.client.HGetAll(ctx, key)
}

func (c *rdc) Del(ctx context.Context, keys ...string) *IntCmd {
	return c.client.Del(ctx, keys...)
}

func (c *rdc) Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd {
	return c.client.Expire(ctx, key, expiration)
}

func (c *rdc) TxPipeline() Pipeliner {
	return c.client.TxPipeline()
}

func (c *rdc) Close() error {
	return c.client.Close()
}
