package kev_test

import (
	"os"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/stretchr/testify/assert"
)

func TestKevProducerTest(t *testing.T) {
	v, ok := os.LookupEnv("TEST_INT")
	if ok && (v == "1" || v == "true" || v == "yes") {
		t.Skip("Skip Integration Test")
	}

	kprod, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
	})

	assert.NoError(t, err)

	defer kprod.Close()

	prod := kev.NewProducer(kprod)

	// ===== KafkaProducer() should return same instance =====
	actual := prod.KafkaProducer()
	if actual != kprod {
		t.Error("KafkaProducer() did not return the original kafka.Producer instance")
	}

	// ===== Events() should return the same event channel =====
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &[]string{"test-topic"}[0], Partition: kafka.PartitionAny},
		Value:          []byte("event test message"),
	}

	deliveryChan := prod.Events()

	err = prod.Produce(msg, deliveryChan)
	if err != nil {
		t.Fatalf("failed to produce message: %v", err)
	}

	select {
	case ev := <-deliveryChan:
		switch m := ev.(type) {
		case *kafka.Message:
			if m.TopicPartition.Error != nil {
				t.Errorf("message delivery failed: %v", m.TopicPartition.Error)
			}
		default:
			t.Errorf("unexpected event type: %T", ev)
		}
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for delivery event from Events()")
	}

	remaining := prod.Flush(5000)
	if remaining > 0 {
		t.Errorf("some messages were not delivered, remaining: %d", remaining)
	}

	prod.Close()
}
