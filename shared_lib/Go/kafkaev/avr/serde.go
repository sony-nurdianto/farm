package avr

import (
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

type avroGenericSerde struct {
	serializer  AvrSerializer
	deserialize AvrDeserializer
}

func NewAvroGenericSerde(client schemaregistry.Client, avr AvrSerdeInstance) (ags avroGenericSerde, err error) {
	ags.serializer, err = avr.NewGenericSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
	if err != nil {
		return ags, err
	}

	ags.deserialize, err = avr.NewGenericDeserializer(client, serde.ValueSerde, avro.NewDeserializerConfig())
	if err != nil {
		return ags, err
	}

	return ags, err
}
