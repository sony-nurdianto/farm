package unit_test

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	mocksClient "github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/test/mocks/schrgs"
	"github.com/stretchr/testify/assert"
)

func TestNewGenericSerializer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocksClient.NewMockClient(ctrl)
	serdeType := serde.ValueSerde
	conf := &avr.SerializerConfig{}

	instance := avr.NewAvrSerdeInstance()
	serializer, err := instance.NewGenericSerializer(mockClient, serdeType, conf)

	assert.NoError(t, err)
	assert.NotNil(t, serializer)
}

func TestNewGenericSerializer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instance := avr.NewAvrSerdeInstance()

	conf := avr.NewSerializerConfig()
	serializer, err := instance.NewGenericSerializer(nil, serde.ValueSerde, conf)

	assert.Error(t, err)
	assert.Nil(t, serializer)
}

func TestNewGenericDeserializer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocksClient.NewMockClient(ctrl)
	serdeType := serde.ValueSerde
	conf := &avro.DeserializerConfig{}

	instance := avr.NewAvrSerdeInstance()
	serializer, err := instance.NewGenericDeserializer(mockClient, serdeType, conf)

	assert.NoError(t, err)
	assert.NotNil(t, serializer)
}

func TestNewGenericDeserializer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	instance := avr.NewAvrSerdeInstance()

	conf := avro.NewDeserializerConfig()
	serializer, err := instance.NewGenericDeserializer(nil, serde.ValueSerde, conf)

	assert.Error(t, err)
	assert.Nil(t, serializer)
}
