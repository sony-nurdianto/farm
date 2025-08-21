package avr

import "github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"

//go:generate mockgen -source=deserializer_def.go -destination=../test/mocks/avr/mock_deserialize.go -package=mocks
type AvrDeserializer interface {
	DeserializeInto(topic string, payload []byte, msg any) error
	Close() error
}

type avrDeserializer struct {
	genericDeserialize *avro.GenericDeserializer
}

func NewAvrDeserializer(gde *avro.GenericDeserializer) *avrDeserializer {
	return &avrDeserializer{
		genericDeserialize: gde,
	}
}

func (s *avrDeserializer) DeserializeInto(topic string, payload []byte, msg any) error {
	return s.genericDeserialize.DeserializeInto(topic, payload, msg)
}

func (s *avrDeserializer) Close() error {
	return s.genericDeserialize.Close()
}
