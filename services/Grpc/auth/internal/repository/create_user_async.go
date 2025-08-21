package repository

import (
	"context"
	"errors"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/constants"
	"github.com/sony-nurdianto/farm/auth/internal/models"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (rp authRepo) publishAvro(
	ctx context.Context,
	accountTopic string,
	userTopic string,
	account *models.InsertAccount,
	user *models.InsertFarmerUser,
) error {
	tracer := otel.Tracer("auth-service")
	_, span := tracer.Start(ctx, "Repo:publishAvro")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", "publish_avro_messages"),
		attribute.String("layer", "repository"),
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.protocol", "avro"),
		attribute.String("messaging.account_topic", accountTopic),
		attribute.String("messaging.user_topic", userTopic),
		attribute.String("user.id", account.Id),
	)

	span.AddEvent("serializing_account_data")
	accountPayload, err := rp.avroSerializer.Serialize(accountTopic, account)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to serialize account data")
		return err
	}
	span.SetAttributes(attribute.Int("messaging.account_payload_size", len(accountPayload)))

	span.AddEvent("serializing_user_data")
	userPayload, err := rp.avroSerializer.Serialize(userTopic, user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to serialize user data")
		return err
	}
	span.SetAttributes(attribute.Int("messaging.user_payload_size", len(userPayload)))

	span.AddEvent("creating_kafka_records")
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

	span.AddEvent("setting_up_transaction_context")
	txCtx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()

	span.AddEvent("beginning_kafka_transaction")
	if err := rp.authProducer.BeginTransaction(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to begin kafka transaction")
		return err
	}

	span.AddEvent("producing_account_message")
	if err := rp.authProducer.Produce(&accountRecord, nil); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to produce account message")
		abortErr := rp.authProducer.AbortTransaction(txCtx)
		if abortErr != nil {
			span.AddEvent("transaction_abort_failed")
			span.RecordError(abortErr)
		}
		return errors.Join(err, abortErr)
	}

	span.AddEvent("producing_user_message")
	if err := rp.authProducer.Produce(&userRecord, nil); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to produce user message")
		abortErr := rp.authProducer.AbortTransaction(txCtx)
		if abortErr != nil {
			span.AddEvent("transaction_abort_failed")
			span.RecordError(abortErr)
		}
		return errors.Join(err, abortErr)
	}

	span.AddEvent("committing_kafka_transaction")
	if err := rp.authProducer.CommitTransaction(txCtx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to commit kafka transaction")
		return err
	}

	span.SetStatus(codes.Ok, "Avro messages published successfully")
	span.SetAttributes(
		attribute.Bool("messaging.transaction_committed", true),
		attribute.Int("messaging.messages_count", 2),
	)

	return nil
}

func (rp authRepo) CreateUserAsync(
	ctx context.Context,
	id, email, fullName, phone, passwordHash string,
) error {
	tracer := otel.Tracer("auth-service")
	cactx, span := tracer.Start(ctx, "Repo:CreateUserAsync")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", "create_user_async"),
		attribute.String("layer", "repository"),
		attribute.String("messaging.operation", "publish"),
		attribute.String("user.id", id),
		attribute.String("user.email", email),
		attribute.String("user.phone", phone),
	)

	span.AddEvent("preparing_data_models")
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

	accountTopic := constants.INSERT_ACCOUNT_TOPIC
	userTopic := constants.INSERT_USER_TOPIC

	span.SetAttributes(
		attribute.String("messaging.account_topic", accountTopic),
		attribute.String("messaging.user_topic", userTopic),
	)

	span.AddEvent("publishing_to_kafka_message_broker")
	err := rp.publishAvro(cactx, accountTopic, userTopic, account, user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to publish messages")
		return err
	}

	span.SetStatus(codes.Ok, "Messages published successfully")
	return nil
}
