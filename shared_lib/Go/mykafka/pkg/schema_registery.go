package pkg

import (
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
)

var SchemaIsNotFoundErr error

type RegisterySchema interface {
	GetLatestSchemaRegistery(subject string) (md schemaregistry.SchemaMetadata, err error)
	GetLatestSchemaID(subject string) (int, error)
	CreateAvroSchema(name string, jsonSchema string, normalize bool) (id int, err error)
	Client() schemaregistry.Client
}

type registerySchema struct {
	client schemaregistry.Client
}

func NewSchemaRegistery(address string, rgs SchemaRegistery) (out registerySchema, _ error) {
	client, err := rgs.NewClient(
		schemaregistry.NewConfig(address),
	)
	if err != nil {
		return out, err
	}

	out.client = client

	return out, nil
}

func (rgs registerySchema) Client() schemaregistry.Client {
	return rgs.client
}

func (rgs registerySchema) GetLatestSchemaRegistery(subject string) (md schemaregistry.SchemaMetadata, err error) {
	md, err = rgs.client.GetLatestSchemaMetadata(subject)
	if err != nil {
		if strings.Contains(err.Error(), "40401") {
			SchemaIsNotFoundErr = err
			return md, SchemaIsNotFoundErr
		}

		return md, err
	}

	return md, nil
}

func (rgs registerySchema) CreateAvroSchema(name string, jsonSchema string, normalize bool) (id int, err error) {
	id, err = rgs.client.Register(name, schemaregistry.SchemaInfo{
		Schema:     jsonSchema,
		SchemaType: "AVRO",
	}, normalize)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (rgs registerySchema) GetLatestSchemaID(subject string) (int, error) {
	schema, err := rgs.client.GetLatestSchemaMetadata(subject)
	if err != nil {
		return -1, err
	}

	return schema.ID, nil
}
