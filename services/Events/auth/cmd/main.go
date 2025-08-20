package main

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
)

func main() {
	// Sync Data User to Redis
	_ = map[kev.ConfigKeyKafka]string{
		kev.BOOTSTRAP_SERVERS:             os.Getenv("KAKFKABROKER"),
		kev.GROUP_ID:                      "auth-consumer-1",
		kev.ENABLE_AUTO_COMMIT:            "true",
		kev.PARTITION_ASSIGNMENT_STRATEGY: "cooperative-sticky",
		kev.AUTO_OFFSET_RESET:             "latest",
	}

	sentinelAddrs := []string{
		"127.0.0.1:26379",
		"127.0.0.1:26380",
		"127.0.0.1:26381",
	}

	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster", // nama master sesuai konfigurasi Sentinel
		SentinelAddrs: sentinelAddrs,
		Password:      "", // jika master punya password
		DB:            0,
	})

	// Test koneksi
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	fmt.Println("Redis connected:", pong)
}
