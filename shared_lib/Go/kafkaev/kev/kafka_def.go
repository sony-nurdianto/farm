package kev

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type KevError = kafka.Error

//go:generate mockgen -source=kafka_def.go -destination=../test/mocks/kev/mock_kafka.go -package=mocks
type Kafka interface {
	NewProducer(conf *kafka.ConfigMap) (KevProducer, error)
	NewConsumer(conf *kafka.ConfigMap) (KevConsumer, error)
}

type kafkaAdapter struct{}

func NewKafka() kafkaAdapter {
	return kafkaAdapter{}
}

func (kafkaAdapter) NewProducer(conf *kafka.ConfigMap) (KevProducer, error) {
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}

	return NewProducer(p), nil
}

func (kafkaAdapter) NewConsumer(conf *kafka.ConfigMap) (KevConsumer, error) {
	c, err := kafka.NewConsumer(conf)
	if err != nil {
		return nil, err
	}

	return NewConsumer(c), nil
}
