package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/repo/insertrepo"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/services/insertsvc"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

func main() {
	godotenv.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	insertFamRepo, err := insertrepo.NewInsertFarmCacheRepo(
		ctx,
		schrgs.NewRegistery(),
		avr.NewAvrSerdeInstance(),
		redis.NewRedisInstance(),
		kev.NewKafka(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	defer insertFamRepo.CloseRepo()

	insertFarmSvc := insertsvc.NewInsertFarmService(insertFamRepo)

	go func(svc insertsvc.InsertFarmCacheService) {
		if err := svc.InsertFarmCache(ctx); err != nil {
			log.Println(err)
		}
	}(insertFarmSvc)

	var once sync.Once

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Event Server Stoping, Gracefully Stop ...")
			fmt.Println("Application Quit.")
			return
		default:
			once.Do(func() { fmt.Println("event service daemon run") })
		}
	}
}
