package pkg

import "github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"

type SchemaRegistery interface {
	NewClient(conf *schemaregistry.Config) (schemaregistry.Client, error)
}

type schemaRegistery struct{}

func NewRegistery() schemaRegistery {
	return schemaRegistery{}
}

func (schemaRegistery) NewClient(conf *schemaregistry.Config) (schemaregistry.Client, error) {
	return schemaregistry.NewClient(conf)
}
