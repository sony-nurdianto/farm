package repository

import (
	"context"
	"errors"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
)

func (rp authRepo) publishAvro(
	accountTopic string,
	userTopic string,
	account *models.InsertAccount,
	user *models.InsertFarmerUser,
) error {
	serializer, err := rp.avro.NewGenericSerializer(
		rp.schemaRegisteryClient.Client(),
		avr.ValueSerde,
		avr.NewSerializerConfig(),
	)
	if err != nil {
		return err
	}

	accountPayload, err := serializer.Serialize(accountTopic, account)
	if err != nil {
		return err
	}

	userPayload, err := serializer.Serialize(userTopic, user)
	if err != nil {
		return err
	}

	accountRecord := kev.MessageKafka{
		TopicPartition: kev.KafkaTopicPartition{
			Topic:     &accountTopic,
			Partition: kev.KafkaPartitionAny,
		},
		Value: accountPayload,
	}.Factory()

	userRecord := kev.MessageKafka{
		TopicPartition: kev.KafkaTopicPartition{
			Topic:     &userTopic,
			Partition: kev.KafkaPartitionAny,
		},
		Value: userPayload,
	}.Factory()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*1000,
	)
	defer cancel()

	if err := rp.authProducer.BeginTransaction(); err != nil {
		return err
	}

	if err := rp.authProducer.Produce(&accountRecord, nil); err != nil {
		return errors.Join(err, rp.authProducer.AbortTransaction(ctx))
	}

	if err := rp.authProducer.Produce(&userRecord, nil); err != nil {
		return errors.Join(err, rp.authProducer.AbortTransaction(ctx))
	}

	return rp.authProducer.CommitTransaction(ctx)
}

func (rp authRepo) CreateUserAsync(
	id, email, fullName, phone, passwordHash string,
) error {
	account := &models.InsertAccount{
		Id:       id,
		Email:    email,
		Password: passwordHash,
	}

	user := &models.InsertFarmerUser{
		Id:       id,
		FullName: fullName,
		Email:    email,
		Phone:    phone,
	}

	accountTopic := "insert-account"
	userTopic := "insert-user"

	return rp.publishAvro(accountTopic, userTopic, account, user)
}
