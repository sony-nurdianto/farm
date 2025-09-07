package insertsvc

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/sony-nurdianto/farm/services/Events/farm/internal/concurrent"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/repo/insertrepo"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
)

type InsertFarmCacheService interface {
	InsertFarmCache(ctx context.Context) error
}

type insertFarmService struct {
	insertFarmRepo insertrepo.InsertFarmCacheRepo
}

func NewInsertFarmService(repo insertrepo.InsertFarmCacheRepo) insertFarmService {
	return insertFarmService{
		insertFarmRepo: repo,
	}
}

func (s insertFarmService) InsertFarmCache(
	ctx context.Context,
) error {
	if err := s.insertFarmRepo.CreateFarmIndex(ctx); err != nil {
		if !errors.Is(err, insertrepo.ErrIndexAlredyExsist) {
			return err
		}
	}

	consumer := s.insertFarmRepo.Consumer()
	if err := consumer.SubscribeTopics([]string{"insert-farm-cache"}, kev.RebalanceCbCooperativeSticky); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if _, ok := err.(kev.KevError); ok {
					continue
				}

				return err
			}

			farms, err := s.insertFarmRepo.FarmDeserializer(msg.Value)
			if err != nil {
				log.Printf("error deserialize farm: %s\n", err.Error())
				continue
			}

			chs := make([]<-chan any, len(farms))

			for i, f := range farms {
				chs[i] = s.insertFarmRepo.InsertFarmCache(ctx, f)
			}

			for v := range concurrent.FanIn(ctx, chs...) {
				res, ok := v.(concurrent.Result[error])
				if !ok {
					log.Println("Wrong Type Data")
					continue
				}

				if res.Error != nil {
					log.Println(res.Error)
					continue
				}
			}

		}
	}
}
