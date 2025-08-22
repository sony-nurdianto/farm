package redis

import "github.com/redis/go-redis/v9"

type RedisInstance interface {
	NewFailoverClient(failoverOpt *FailoverOptions) *Client
}

type rdi struct{}

func NewRedisInstance() rdi {
	return rdi{}
}

func (rdi) NewFailoverClient(failoverOpt *FailoverOptions) RedisClient {
	rdc := redis.NewFailoverClient(failoverOpt)

	return NewRedisClient(rdc)
}
