package unit_test

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	mocks "github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/test/mocks/avr"
	"github.com/stretchr/testify/assert"
)

func TestAvrDeserializer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeserializer := mocks.NewMockAvrDeserializer(ctrl)
	topic := "test-topic"
	payload := []byte("mock-avro-data")
	var msg DummyStruct

	mockDeserializer.EXPECT().
		DeserializeInto(topic, payload, &msg).
		DoAndReturn(func(topic string, payload []byte, msg any) error {
			// Simulate deserialization
			if dummyMsg, ok := msg.(*DummyStruct); ok {
				dummyMsg.Name = "Deserialized Name"
			}
			return nil
		})

	// Test
	err := mockDeserializer.DeserializeInto(topic, payload, &msg)
	assert.NoError(t, err)
	assert.Equal(t, "Deserialized Name", msg.Name)
}

func TestAvrDeserializer_DeserializeInto_Implementation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGenericDeserializer := &avro.GenericDeserializer{}
	deserializer := avr.NewAvrDeserializer(mockGenericDeserializer)

	assert.NotNil(t, deserializer)

	topic := "test-topic"
	payload := []byte("test-payload")
	var msg DummyStruct

	err := deserializer.DeserializeInto(topic, payload, &msg)

	assert.Error(t, err)
}
