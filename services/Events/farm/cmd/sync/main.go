package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/vfs"
	"github.com/joho/godotenv"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/repo/syncrepo"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/services/syncsvc"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

func main() {
	godotenv.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	opts := &pebble.Options{
		FS: vfs.Default,
	}

	stateDB, err := pebble.Open("state", opts)
	if err != nil {
		log.Fatalln(err)
	}

	farmRepo, err := syncrepo.NewSyncFarmRepo(
		ctx,
		schrgs.NewRegistery(),
		avr.NewAvrSerdeInstance(),
		kev.NewKafka(),
		redis.NewRedisInstance(),
		stateDB,
	)
	if err != nil {
		log.Fatalln(err)
	}

	defer farmRepo.CloseRepo()

	farmService := syncsvc.NewFarmService(farmRepo)

	go func(c context.Context, fs syncsvc.FarmService) {
		if err := fs.SyncFarmCache(c, "farm-db.public.farms_all_partitions"); err != nil {
			log.Fatalln(err)
		}
	}(ctx, farmService)

	go func(c context.Context, fs syncsvc.FarmService) {
		if err := fs.SyncFarmAddressCache(c, "farm-db.public.addresses_all_partitions"); err != nil {
			log.Fatalln(err)
		}
	}(ctx, farmService)

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
