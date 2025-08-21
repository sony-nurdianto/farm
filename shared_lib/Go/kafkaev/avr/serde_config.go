package avr

import (
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

type SerializerConfig struct {
	avro.SerializerConfig
}

func NewSerializerConfig() *SerializerConfig {
	conf := avro.NewSerializerConfig()

	return &SerializerConfig{
		SerializerConfig: *conf,
	}
}

func (c *SerializerConfig) ToAvroConfig() *avro.SerializerConfig {
	return &c.SerializerConfig
}

type DeserializerConfig struct {
	avro.DeserializerConfig
}

func NewDeserializerConfig() *DeserializerConfig {
	conf := avro.NewDeserializerConfig()

	return &DeserializerConfig{
		DeserializerConfig: *conf,
	}
}

func (c *DeserializerConfig) ToAvroConfig() *avro.DeserializerConfig {
	return &c.DeserializerConfig
}
