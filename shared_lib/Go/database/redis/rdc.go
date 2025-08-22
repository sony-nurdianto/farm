package redis

import (
	"context"
)

type RedisClient interface {
	HSet(ctx context.Context, key string, values ...any) *IntCmd
	HGet(ctx context.Context, key string, field string) *StringCmd
	HGetAll(ctx context.Context, key string) *MapStringStringCmd
	Ping(ctx context.Context) StatusCmd
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
