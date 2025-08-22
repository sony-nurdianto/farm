package redis

import "github.com/redis/go-redis/v9"

type (
	FailoverOptions    = redis.FailoverOptions
	Client             = redis.Client
	IntCmd             = redis.IntCmd
	StringCmd          = redis.StringCmd
	MapStringStringCmd = redis.MapStringStringCmd
)
