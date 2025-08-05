package pkg

import (
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

type avroGenericSerde struct {
	serializer  *avro.GenericSerializer
	deserialize *avro.GenericDeserializer
}

func NewAvroGenericSerde(client schemaregistry.Client) (ags avroGenericSerde, err error) {
	ags.serializer, err = avro.NewGenericSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
	if err != nil {
		return ags, err
	}

	ags.deserialize, err = avro.NewGenericDeserializer(client, serde.ValueSerde, avro.NewDeserializerConfig())
	if err != nil {
		return ags, err
	}

	return ags, err
}

func (ags avroGenericSerde) Serialize(topic string, msg any) ([]byte, error) {
	return ags.serializer.Serialize(topic, msg)
}

func (ags avroGenericSerde) DeserializeInto(topic string, payload []byte, msg any) error {
	return ags.deserialize.DeserializeInto(topic, payload, msg)
}
