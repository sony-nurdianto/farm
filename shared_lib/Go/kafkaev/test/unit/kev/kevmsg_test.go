package unit_test

import (
	"testing"

	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/stretchr/testify/assert"
)

func TestKevMessage(t *testing.T) {
	t.Run("Header", func(t *testing.T) {
		h := kev.KafkaHeader{
			Key:   "Something",
			Value: []byte("Value"),
		}

		out := h.Factory()
		assert.Equal(t, h.Key, out.Key)
		assert.Equal(t, h.Value, out.Value)
	})

	t.Run("Message", func(t *testing.T) {
		header := kev.KafkaHeader{
			Key:   "Something",
			Value: []byte("Value"),
		}

		tpc := "topic"
		kevMsg := kev.MessageKafka{
			TopicPartition: kev.KafkaTopicPartition{
				Topic:     &tpc,
				Partition: kev.KafkaPartitionAny,
			},
			Headers: []kev.KafkaHeader{header},
			Key:     []byte("Something"),
			Value:   []byte("Something"),
		}

		factory := kevMsg.Factory()
		assert.Equal(t, factory.Headers[0].Key, kevMsg.Headers[0].Key)
	})
}
