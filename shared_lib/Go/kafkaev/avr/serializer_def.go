package avr

import "github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"

//go:generate mockgen -source=serializer_def.go -destination=../test/mocks/avr/mock_serializer.go -package=mocks
type AvrSerializer interface {
	Serialize(topic string, msg any) ([]byte, error)
}

type avrSerializer struct {
	genericSerializer *avro.GenericSerializer
}

func NewAvrSerializer(gse *avro.GenericSerializer) *avrSerializer {
	return &avrSerializer{
		genericSerializer: gse,
	}
}

func (s *avrSerializer) Serialize(topic string, msg any) ([]byte, error) {
	return s.genericSerializer.Serialize(topic, msg)
}
