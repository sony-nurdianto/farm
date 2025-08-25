package repo

import (
	"context"
	"time"

	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/concurent"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/constants"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
)

func (fr farmerRepo) publishAccountAvro(
	ctx context.Context,
	accountTopic string,
	account *models.Accounts,
) <-chan any {

	out := make(chan any, 1)

	go func() {
		defer close(out)

		payload, err := fr.avrSerializer.Serialize(accountTopic, account)
		if err != nil {
			send(ctx, out, err)
			return
		}

		record := kev.MessageKafka{
			TopicPartition: kev.KafkaTopicPartition{
				Topic:     &accountTopic,
				Partition: kev.KafkaPartitionAny,
			},
			Value: payload,
		}.Factory()

		if err := fr.farmerProudcer.Produce(&record, nil); err != nil {
			send(ctx, out, err)
			return
		}

		send(ctx, out, nil)
	}()

	return out

}

func (fr farmerRepo) publishUserAvro(
	ctx context.Context,
	userTopic string,
	users *models.UpdateUsers,
) <-chan any {
	out := make(chan any, 1)

	go func() {
		defer close(out)

		payload, err := fr.avrSerializer.Serialize(userTopic, users)
		if err != nil {
			send(ctx, out, err)
		}

		record := kev.MessageKafka{
			TopicPartition: kev.KafkaTopicPartition{
				Topic:     &userTopic,
				Partition: kev.KafkaPartitionAny,
			},
			Value: payload,
		}.Factory()

		if err := fr.farmerProudcer.Produce(&record, nil); err != nil {
			send(ctx, out, err)
			return
		}

		send(ctx, out, nil)
	}()

	return out
}

func (fr farmerRepo) UpdateUserAsync(ctx context.Context, users *models.UpdateUsers) error {

	fCtx, done := context.WithTimeout(ctx, time.Second*15)
	defer done()

	chs := make([]<-chan any, 0)
	if users.Email != nil {
		publishAccount := fr.publishAccountAvro(
			fCtx,
			constants.UpdateAccountTopic,
			&models.Accounts{
				ID:    users.ID,
				Email: *users.Email,
			},
		)

		chs = append(chs, publishAccount)
	}

	publishUser := fr.publishUserAvro(ctx, constants.UpdateUsersTopic, users)
	chs = append(chs, publishUser)

	for res := range concurent.FanIn(fCtx, chs...) {
		switch v := res.(type) {
		case error:
			return v

		}
	}

	return nil
}
