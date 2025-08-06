package schrgs

import "github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"

//go:generate mockgen -source=schemargs_instance_def.go -destination=../test/mocks/schrgs/mock_schemargs_instance.go -package=mocks
type SchemaRegisteryInstance interface {
	NewClient(conf *schemaregistry.Config) (schemaregistry.Client, error)
}

type schemaRegisteryInstance struct{}

func NewRegistery() schemaRegisteryInstance {
	return schemaRegisteryInstance{}
}

func (schemaRegisteryInstance) NewClient(conf *schemaregistry.Config) (schemaregistry.Client, error) {
	return schemaregistry.NewClient(conf)
}
