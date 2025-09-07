package syncsvc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sony-nurdianto/farm/services/Events/farm/internal/repo/syncrepo"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
)

type FarmService interface {
	SyncFarmCache(ctx context.Context, topic string) error
	SyncFarmAddressCache(ctx context.Context, topic string) error
}

type farmService struct {
	syncRepo syncrepo.SyncFarmRepo
}

func NewFarmService(repo syncrepo.SyncFarmRepo) farmService {
	return farmService{
		repo,
	}
}

func (fs farmService) SyncFarmAddressCache(
	ctx context.Context,
	topic string,
) error {
	consumer := fs.syncRepo.FarmAddrConsumer()
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

			farmAddr, err := fs.syncRepo.DeserializerFarmAddress(topic, msg.Value)
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
				if err := fs.syncRepo.UpsertFarmAddressCache(ctx, farmAddr); err != nil {
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
	consumer := fs.syncRepo.FarmConsumer()
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

			farm, err := fs.syncRepo.DeserializerFarm(topic, msg.Value)
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
				if err := fs.syncRepo.UpsertFarmCache(ctx, farm, op); err != nil {
					log.Println(err)
					continue
				}
			case "d":
				if err := fs.syncRepo.DeleteFarmCache(ctx, farm.ID, farm.FarmerID); err != nil {
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
