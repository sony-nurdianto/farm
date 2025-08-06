package kev

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

//go:generate mockgen -source=kafka_def.go -destination=../test/mocks/kev/mock_kafka.go -package=mocks
type Kafka interface {
	NewProducer(conf *kafka.ConfigMap) (KevProducer, error)
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
