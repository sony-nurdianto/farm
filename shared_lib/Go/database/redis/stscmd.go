package redis

import "github.com/redis/go-redis/v9"

type StatusCmd interface {
	Result() (string, error)
}

type stsCmd struct {
	sts *redis.StatusCmd
}

func NewStatusCmd(sts *redis.StatusCmd) *stsCmd {
	return &stsCmd{sts}
}

func (s *stsCmd) Result() (string, error) {
	return s.sts.Result()
}
