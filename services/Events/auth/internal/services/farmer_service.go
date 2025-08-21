package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sony-nurdianto/farm/services/Events/auth/internal/repo"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
)

type farmerService struct {
	repo repo.Repo
}

func NewFarmerService(rp repo.Repo) farmerService {
	return farmerService{
		repo: rp,
	}
}

func (fs farmerService) SyncUserCache(ctx context.Context, topic string) error {
	consumer := fs.repo.Consumer()
	consumer.SubscribeTopics([]string{topic}, kev.RebalanceCbCooperativeSticky)

	for {
		msg, err := consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			if _, ok := err.(kev.KevError); ok {
				continue
			}

			return err
		}

		farmer, err := fs.repo.DeserializerFarmer(topic, msg.Value)
		if err != nil {
			log.Println(err)
		}

		cacheKey := fmt.Sprintf("farmer:%s", farmer.ID)

		if err := fs.repo.InsertFarmerCache(ctx, cacheKey, farmer); err != nil {
			log.Println(err)
		}

		if _, err := consumer.CommitMessage(msg); err != nil {
			return err
		}

	}
}
