package pkg

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type Kafka interface {
	NewProducer(conf *kafka.ConfigMap) (Producer, error)
}

type kafkaAdapter struct{}

func NewKafka() kafkaAdapter {
	return kafkaAdapter{}
}

func (kafkaAdapter) NewProducer(conf *kafka.ConfigMap) (Producer, error) {
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, err
	}

	return NewProducer(p), nil
}
