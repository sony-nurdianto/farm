package schrgs

import (
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
)

var SchemaIsNotFoundErr error

type SchemaRegistery struct {
	client SchemaRegisteryClient
}

func NewSchemaRegistery(address string, rgi SchemaRegisteryInstance) (out SchemaRegistery, _ error) {
	client, err := rgi.NewClient(
		schemaregistry.NewConfig(address),
	)
	if err != nil {
		return out, err
	}

	out.client = NewSchemaRegisteryClient(client)

	return out, nil
}

func (rgs SchemaRegistery) Client() SchemaRegisteryClient {
	return rgs.client
}

func (rgs SchemaRegistery) GetLatestSchemaRegistery(subject string) (md schemaregistry.SchemaMetadata, err error) {
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

func (rgs SchemaRegistery) CreateAvroSchema(name string, jsonSchema string, normalize bool) (id int, err error) {
	id, err = rgs.client.Register(name, schemaregistry.SchemaInfo{
		Schema:     jsonSchema,
		SchemaType: "AVRO",
	}, normalize)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (rgs SchemaRegistery) GetLatestSchemaID(subject string) (int, error) {
	schema, err := rgs.client.GetLatestSchemaMetadata(subject)
	if err != nil {
		return -1, err
	}

	return schema.ID, nil
}
