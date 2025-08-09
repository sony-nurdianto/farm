package avr

import (
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

type avroGenericSerde struct {
	Serializer  AvrSerializer
	Deserialize AvrDeserializer
}

func NewAvroGenericSerde(client schemaregistry.Client, avr AvrSerdeInstance) (ags avroGenericSerde, err error) {
	ags.Serializer, err = avr.NewGenericSerializer(client, serde.ValueSerde, NewSerializerConfig())
	if err != nil {
		return ags, err
	}

	ags.Deserialize, err = avr.NewGenericDeserializer(client, serde.ValueSerde, avro.NewDeserializerConfig())
	if err != nil {
		return ags, err
	}

	return ags, err
}
