package schrgs

import "github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"

type SchemaRegisteryConfig = schemaregistry.Config

//go:generate mockgen -source=schemargs_instance_def.go -destination=../test/mocks/schrgs/mock_schemargs_instance.go -package=mocks
type SchemaRegisteryInstance interface {
	NewClient(conf *schemaregistry.Config) (schemaregistry.Client, error)
	NewConfig(url string) *schemaregistry.Config
}

type schemaRegisteryInstance struct{}

func NewRegistery() schemaRegisteryInstance {
	return schemaRegisteryInstance{}
}

func (schemaRegisteryInstance) NewClient(conf *schemaregistry.Config) (schemaregistry.Client, error) {
	return schemaregistry.NewClient(conf)
}

func (schemaRegisteryInstance) NewConfig(url string) *schemaregistry.Config {
	return schemaregistry.NewConfig(url)
}
