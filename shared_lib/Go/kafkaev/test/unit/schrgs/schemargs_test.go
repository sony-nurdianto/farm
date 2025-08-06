package schrgs_test

import (
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
	mocks "github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/test/mocks/schrgs"
	"github.com/stretchr/testify/assert"
)

func TestNewSchemaRegistery_NewClient_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInstance := mocks.NewMockSchemaRegisteryInstance(ctrl)
	mockInstance.EXPECT().NewClient(gomock.Any()).Return(nil, errors.New("Failed To Create NewClient"))

	rgsSchema, err := schrgs.NewSchemaRegistery("someaddress", mockInstance)
	assert.Error(t, err)
	assert.Empty(t, rgsSchema)
}

func TestSchemaRegistery_GetLatestSchemaRegistery_Error(t *testing.T) {
	t.Run("Error SchemaIsNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKafkaClient := mocks.NewMockClient(ctrl) // schemaregistry.Client
		mockInstance := mocks.NewMockSchemaRegisteryInstance(ctrl)

		mockKafkaClient.EXPECT().
			GetLatestSchemaMetadata("something-subject").
			Return(schemaregistry.SchemaMetadata{}, errors.New("40401: subject not found"))

		mockInstance.EXPECT().
			NewClient(gomock.Any()).
			Return(mockKafkaClient, nil)

		rgsSchema, err := schrgs.NewSchemaRegistery("something", mockInstance)
		assert.NoError(t, err)

		schrgs.SchemaIsNotFoundErr = nil

		md, err := rgsSchema.GetLatestSchemaRegistery("something-subject")

		assert.Error(t, err)
		assert.Equal(t, schrgs.SchemaIsNotFoundErr, err)
		assert.Equal(t, schemaregistry.SchemaMetadata{}, md)
	})
}
