package services

import (
	"context"
	"strings"
	"time"

	"github.com/sony-nurdianto/farm/services/Events/farm/internal/models"
	"github.com/sony-nurdianto/farm/services/Events/farm/internal/repo"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
)

type FarmService interface {
	SyncFarmCache(
		ctx context.Context,
		topic string,
	) error
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
	consumer := fs.repo.FarmConsumer()
	consumer.SubscribeTopics([]string{topic}, kev.RebalanceCbCooperativeSticky)

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

			var farmAddr models.FarmAddress

			err = fs.repo.DeserializerFarm(topic, msg.Value, &farmAddr)
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

		}
	}
}

func (fs farmService) SyncFarmCache(
	ctx context.Context,
	topic string,
) error {
	consumer := fs.repo.FarmConsumer()
	consumer.SubscribeTopics([]string{topic}, kev.RebalanceCbCooperativeSticky)

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

			var farm models.Farm

			err = fs.repo.DeserializerFarm(topic, msg.Value, &farm)
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
				if err := fs.repo.UpsertFarmCache(ctx, farm, op); err != nil {
					continue
				}
			case "d":
				if err := fs.repo.DeleteFarmCache(ctx, farm.ID, farm.AddressesID, farm.FarmerID); err != nil {
					continue
				}
			}

		}
	}
}
