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
	FTCreate(ctx context.Context, index string, options *FTCreateOptions, schema ...*FieldSchema) StatusCmd
	FTSearch(ctx context.Context, index string, query string) *FTSearchCmd
	FTSearchWithArgs(ctx context.Context, index string, query string, options *FTSearchOptions) *FTSearchCmd
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

func (c *rdc) Do(ctx context.Context, args ...any) *RedisCmd {
	return c.client.Do(ctx, args...)
}

func (c *rdc) Expire(ctx context.Context, key string, expiration time.Duration) *BoolCmd {
	return c.client.Expire(ctx, key, expiration)
}

func (c *rdc) TxPipeline() Pipeliner {
	return c.client.TxPipeline()
}

func (c *rdc) FTCreate(ctx context.Context, index string, options *FTCreateOptions, schema ...*FieldSchema) StatusCmd {
	return c.client.FTCreate(ctx, index, options, schema...)
}

func (c *rdc) FTSearch(ctx context.Context, index string, query string) *FTSearchCmd {
	return c.client.FTSearch(ctx, index, query)
}

func (c *rdc) FTSearchWithArgs(ctx context.Context, index string, query string, options *FTSearchOptions) *FTSearchCmd {
	return c.client.FTSearchWithArgs(ctx, index, query, options)
}

func (c *rdc) Close() error {
	return c.client.Close()
}
