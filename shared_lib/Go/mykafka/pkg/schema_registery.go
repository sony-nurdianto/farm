package pkg

import (
	"strings"

	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
)

var SchemaIsNotFoundErr error

type SchemaRegistery interface {
	GetSchemaRegistery(subject string, versrion int) (md schemaregistry.SchemaMetadata, err error)
	RegisterSchema(name string, jsonSchema string, normalize bool) (id int, err error)
	Client() schemaregistry.Client
}

type registerySchema struct {
	client schemaregistry.Client
}

func NewSchemaRegistery(address string) (rgs registerySchema, _ error) {
	client, err := schemaregistry.NewClient(
		schemaregistry.NewConfig(address),
	)
	if err != nil {
		return rgs, err
	}

	rgs.client = client

	return rgs, nil
}

func (rgs registerySchema) Client() schemaregistry.Client {
	return rgs.client
}

func (rgs registerySchema) GetSchemaRegistery(subject string, versrion int) (md schemaregistry.SchemaMetadata, err error) {
	md, err = rgs.client.GetSchemaMetadata(subject, versrion)
	if err != nil {
		if strings.Contains(err.Error(), "40401") {
			SchemaIsNotFoundErr = err
			return md, SchemaIsNotFoundErr
		}

		return md, err
	}

	return md, nil
}

func (rgs registerySchema) RegisterSchema(name string, jsonSchema string, normalize bool) (id int, err error) {
	id, err = rgs.client.Register(name, schemaregistry.SchemaInfo{
		Schema:     jsonSchema,
		SchemaType: "AVRO",
	}, normalize)
	if err != nil {
		return 0, err
	}

	return id, err
}
