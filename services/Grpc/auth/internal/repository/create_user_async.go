package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
)

func (rp authRepo) ensureSchema(
	subject string,
	schema string,
) error {
	_, err := rp.schemaRegistery.GetLatestSchemaRegistery(subject)
	if err == nil {
		return nil
	}

	if !errors.Is(err, schrgs.SchemaIsNotFoundErr) {
		return err
	}

	_, regErr := rp.schemaRegistery.CreateAvroSchema(subject, schema, false)
	if regErr != nil {
		return regErr
	}

	return nil
}

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
		// avro.NewSerializerConfig(),
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

	cfg := map[kev.ConfigKeyKafka]string{
		kev.BOOTSTRAP_SERVERS:  "localhost:29092",
		kev.ACKS:               "all",
		kev.ENABLE_IDEMPOTENCE: "true",
		kev.COMPRESION_TYPE:    "snappy",
		kev.RETRIES:            "5",
		kev.RETRY_BACKOFF_MS:   "100",
		kev.LINGER_MS:          "5",
		kev.MAX_IN_FLIGHT_REQUESTS_PER_CONNECTION: "5",

		kev.TRANSACTIONAL_ID: "register-user",
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

	producer, err := rp.kafka.Producer(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond*1000,
	)
	defer cancel()

	// for i := range 5 {
	// 	err := producer.InitTransactions(ctx)
	// 	if err == nil {
	// 		fmt.Printf("✅ Kafka transactional producer ready after %d attempt(s)\n", i)
	// 		break
	// 	}
	// 	log.Println("Waiting Init Transaction Ready")
	// 	time.Sleep(2 * time.Second)
	// }

	if err := producer.InitTransactions(ctx); err != nil {
		return err
	}

	if err := producer.BeginTransaction(); err != nil {
		return err
	}

	if err := producer.Produce(&accountRecord, nil); err != nil {
		return errors.Join(err, producer.AbortTransaction(ctx))
	}

	if err := producer.Produce(&userRecord, nil); err != nil {
		return errors.Join(err, producer.AbortTransaction(ctx))
	}

	return producer.CommitTransaction(ctx)
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

	valueStr := "-value"
	accountTopic := "insert-account"
	accountSubject := fmt.Sprintf("%s-%s", accountTopic, valueStr)

	userTopic := "insert-user"
	userSubject := fmt.Sprintf("%s-%s", userTopic, valueStr)

	err := rp.ensureSchema(accountSubject, account.Schema())
	if err != nil {
		return err
	}

	err = rp.ensureSchema(userSubject, user.Schema())
	if err != nil {
		return err
	}

	return rp.publishAvro(accountTopic, userTopic, account, user)
}
