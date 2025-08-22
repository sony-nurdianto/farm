package schrgs

import "github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"

type SchrgsClient = schemaregistry.Client

//go:generate mockgen -source=schemargs_client_def.go -destination=../test/mocks/schrgs/mock_schemargs_client.go -package=mocks
//go:generate mockgen -destination=../test/mocks/schrgs/mock_confluent_client.go -package=mocks github.com/confluentinc/confluent-kafka-go/v2/schemaregistry Client
type SchemaRegisteryClient interface {
	GetLatestSchemaMetadata(subject string) (schemaregistry.SchemaMetadata, error)
	Register(subject string, schema schemaregistry.SchemaInfo, normalize bool) (id int, err error)
	Client() schemaregistry.Client
}

type schemaRegisteryClient struct {
	client schemaregistry.Client
}

func NewSchemaRegisteryClient(client schemaregistry.Client) schemaRegisteryClient {
	return schemaRegisteryClient{client}
}

func (c schemaRegisteryClient) GetLatestSchemaMetadata(subject string) (schemaregistry.SchemaMetadata, error) {
	return c.client.GetLatestSchemaMetadata(subject)
}

func (c schemaRegisteryClient) Register(subject string, schema schemaregistry.SchemaInfo, normalize bool) (id int, err error) {
	return c.client.Register(subject, schema, normalize)
}

func (c schemaRegisteryClient) Client() schemaregistry.Client {
	return c.client
}
