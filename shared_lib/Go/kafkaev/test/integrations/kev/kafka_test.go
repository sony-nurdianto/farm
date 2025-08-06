package kev_test

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
)

func TestKafkaAdapter_NewProducer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conf := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}

	kafkaAdapter := kev.NewKafka()
	producer, err := kafkaAdapter.NewProducer(conf)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if producer == nil {
		t.Error("Expected producer, got nil")
	}
}
