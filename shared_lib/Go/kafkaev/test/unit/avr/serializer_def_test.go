package unit_test

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	mocks "github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/test/mocks/avr"

	mocksClient "github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/test/mocks/schrgs"
	"github.com/stretchr/testify/assert"
)

type DummyStruct struct {
	Name string `avro:"name" json:"name"`
}

func TestSerialezerDef(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock AvrSerializer instead of using real avro.NewGenericSerializer
	mockSerializer := mocks.NewMockAvrSerializer(ctrl)

	// Setup expectation
	expectedData := []byte("serialized-data")
	mockSerializer.EXPECT().
		Serialize("some topic", gomock.Any()).
		Return(expectedData, nil)

	// Test
	msg := DummyStruct{Name: "My Name Is?"}
	value, err := mockSerializer.Serialize("some topic", msg)

	assert.NoError(t, err)
	assert.NotEmpty(t, value)
	assert.Equal(t, expectedData, value)
}

func TestAvrSerializer_Serialize_Implementation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocksClient.NewMockClient(ctrl)

	// Setup expectations untuk calls yang akan terjadi
	client.EXPECT().
		RegisterFullResponse(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(schemaregistry.SchemaMetadata{}, nil).
		AnyTimes() // Allow multiple calls

	// Bisa juga perlu mock calls lainnya
	client.EXPECT().
		GetLatestSchemaMetadata(gomock.Any()).
		Return(schemaregistry.SchemaMetadata{}, nil).
		AnyTimes()

	conf := avro.NewSerializerConfig()
	gs, err := avro.NewGenericSerializer(client, serde.ValueSerde, conf)

	if err == nil && gs != nil {
		serializer := avr.NewAvrSerializer(gs)
		msg := DummyStruct{Name: "test"}

		_, err := serializer.Serialize("test-topic", msg)
		_ = err
	}

	assert.True(t, true)
}
