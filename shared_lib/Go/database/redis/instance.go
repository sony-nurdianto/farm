package redis

import "github.com/redis/go-redis/v9"

type RedisInstance interface {
	NewFailoverClient(failoverOpt *FailoverOptions) *Client
}

type rdi struct{}

func NewRedisInstance() rdi {
	return rdi{}
}

func (rdi) NewFailoverClient(failoverOpt *FailoverOptions) *Client {
	return redis.NewFailoverClient(failoverOpt)
}
