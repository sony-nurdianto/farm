package schrgs

import (
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
)

var SchemaIsNotFoundErr error

type registerySchema struct {
	client SchemaRegisteryClient
}

func NewSchemaRegistery(address string, rgi SchemaRegisteryInstance) (out registerySchema, _ error) {
	client, err := rgi.NewClient(
		schemaregistry.NewConfig(address),
	)
	if err != nil {
		return out, err
	}

	out.client = NewSchemaRegisteryClient(client)

	return out, nil
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
