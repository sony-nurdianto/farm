package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sony-nurdianto/farm/services/Events/auth/internal/repo"
	"github.com/sony-nurdianto/farm/services/Events/auth/internal/services"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

func redisConnection() (c *redis.Client, err error) {
	count := 0
	for range 5 {
		rdb := redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    os.Getenv("FARMER_REDIS_MASTER_NAME"),
			SentinelAddrs: []string{os.Getenv("SENTINEL_FARMER_REDIS_ADDR")},
			Username:      os.Getenv("FARMER_REDIS_MASTER_USER_NAME"),
			Password:      os.Getenv("FARMER_REDIS_MASTER_PASSWORD"),
			DB:            0,
		})

		_, err = rdb.Ping(context.Background()).Result()
		if err == nil {
			return rdb, nil
		}

		time.Sleep(time.Second * 2)
		count++
	}

	return nil, fmt.Errorf("connect failed after %d attempts: %w", count, err)
}

func main() {
	godotenv.Load()
	rdb, err := redisConnection()
	if err != nil {
		log.Fatalln(err)
	}

	cfgFarmer := map[kev.ConfigKeyKafka]string{
		kev.BOOTSTRAP_SERVERS:             os.Getenv("KAKFKABROKER"),
		kev.GROUP_ID:                      "auth-event",
		kev.AUTO_OFFSET_RESET:             "earliest",
		kev.ENABLE_AUTO_COMMIT:            "false",
		kev.PARTITION_ASSIGNMENT_STRATEGY: "cooperative-sticky",
	}

	farmerSvcRepo, err := repo.NewAuthRepo(
		schrgs.NewRegistery(),
		avr.NewAvrSerdeInstance(),
		kev.NewKafka(),
		cfgFarmer,
		rdb,
	)
	if err != nil {
		log.Fatalln(err)
	}

	defer farmerSvcRepo.CloseRepo()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		farmerSvc := services.NewFarmerService(farmerSvcRepo)
		if err := farmerSvc.SyncUserCache(ctx, "farmer-db.public.users_all_partitions"); err != nil {
			log.Fatalln(err)
		}
	}()

	fmt.Println("event service daemon run")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Event Server Stoping, Gracefully Stop ...")
			fmt.Println("Application Quit.")
			return
		default:
		}
	}
}
