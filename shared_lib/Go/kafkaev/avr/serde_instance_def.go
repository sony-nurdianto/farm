package avr

import (
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

//go:generate mockgen -source=serde_instance_def.go -destination=../test/mocks/avr/mock_serde_instance.go -package=mocks
type AvrSerdeInstance interface {
	NewGenericSerializer(
		client schemaregistry.Client,
		serdeType serde.Type,
		conf *SerializerConfig,
	) (AvrSerializer, error)

	NewGenericDeserializer(
		client schemaregistry.Client,
		serdeType serde.Type,
		conf *DeserializerConfig,
	) (AvrDeserializer, error)
}

type avrSerdeInstance struct{}

func NewAvrSerdeInstance() avrSerdeInstance {
	return avrSerdeInstance{}
}

func (asi avrSerdeInstance) NewGenericSerializer(
	client schemaregistry.Client,
	serdeType serde.Type,
	conf *SerializerConfig,
) (AvrSerializer, error) {
	gs, err := avro.NewGenericSerializer(client, serdeType, conf.ToAvroConfig())
	if err != nil {
		return nil, err
	}

	return NewAvrSerializer(gs), err
}

func (asi avrSerdeInstance) NewGenericDeserializer(
	client schemaregistry.Client,
	serdeType serde.Type,
	conf *DeserializerConfig,
) (AvrDeserializer, error) {
	gds, err := avro.NewGenericDeserializer(client, serdeType, conf.ToAvroConfig())
	if err != nil {
		return nil, err
	}

	return NewAvrDeserializer(gds), nil
}
