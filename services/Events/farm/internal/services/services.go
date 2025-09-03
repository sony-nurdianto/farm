package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sony-nurdianto/farm/services/Events/farm/internal/repo"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
)

type FarmService interface {
	SyncFarmCache(ctx context.Context, topic string) error
	SyncFarmAddressCache(ctx context.Context, topic string) error
}

type farmService struct {
	repo repo.FarmRepo
}

func NewFarmService(repo repo.FarmRepo) farmService {
	return farmService{
		repo,
	}
}

func (fs farmService) SyncFarmAddressCache(
	ctx context.Context,
	topic string,
) error {
	consumer := fs.repo.FarmAddrConsumer()
	err := consumer.SubscribeTopics([]string{topic}, kev.RebalanceCbCooperativeSticky)
	if err != nil {
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

			farmAddr, err := fs.repo.DeserializerFarmAddress(topic, msg.Value)
			if err != nil {
				continue
			}

			var op string
			for _, h := range msg.Headers {
				key := strings.TrimPrefix(h.Key, "__")
				if key == "op" {
					op = string(h.Value)
					break
				}
			}

			switch op {
			case "c", "u", "r":
				if err := fs.repo.UpsertFarmAddressCache(ctx, farmAddr); err != nil {
					continue
				}
			case "d":
				continue
			}

			if _, err := consumer.CommitMessage(msg); err != nil {
				fmt.Printf("Error committing message: %v\n", err)
				continue
			}

		}
	}
}

func (fs farmService) SyncFarmCache(
	ctx context.Context,
	topic string,
) error {
	consumer := fs.repo.FarmConsumer()
	if err := consumer.SubscribeTopics([]string{topic}, kev.RebalanceCbCooperativeSticky); err != nil {
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

			farm, err := fs.repo.DeserializerFarm(topic, msg.Value)
			if err != nil {
				log.Println(err)
				continue
			}

			var op string
			for _, h := range msg.Headers {
				key := strings.TrimPrefix(h.Key, "__")
				if key == "op" {
					op = string(h.Value)
					break
				}
			}

			switch op {
			case "c", "u", "r":
				if err := fs.repo.UpsertFarmCache(ctx, farm, op); err != nil {
					log.Println(err)
					continue
				}
			case "d":
				if err := fs.repo.DeleteFarmCache(ctx, farm.ID, farm.AddressID, farm.FarmerID); err != nil {
					log.Println(err)
					continue
				}
			}

			if _, err := consumer.CommitMessage(msg); err != nil {
				fmt.Printf("Error committing message: %v\n", err)
				continue
			}
		}
	}
}
