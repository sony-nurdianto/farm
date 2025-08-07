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

func TestSchemaRegistery_GetLatestSchemaRegistery_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockKafkaClient := mocks.NewMockClient(ctrl)
	mockInstance := mocks.NewMockSchemaRegisteryInstance(ctrl)

	mockKafkaClient.EXPECT().
		GetLatestSchemaMetadata(gomock.Any()).
		Return(schemaregistry.SchemaMetadata{}, nil)

	mockInstance.EXPECT().
		NewClient(gomock.Any()).
		Return(mockKafkaClient, nil)

	rgsSchema, err := schrgs.NewSchemaRegistery(gomock.Any().String(), mockInstance)
	assert.NoError(t, err)

	md, err := rgsSchema.GetLatestSchemaRegistery("something")
	assert.NoError(t, err)

	assert.Equal(t, md, schemaregistry.SchemaMetadata{})
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

	t.Run("Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockKafkaClient := mocks.NewMockClient(ctrl) // schemaregistry.Client
		mockInstance := mocks.NewMockSchemaRegisteryInstance(ctrl)

		mockKafkaClient.EXPECT().
			GetLatestSchemaMetadata("something-subject").
			Return(schemaregistry.SchemaMetadata{}, errors.New("Other Error"))

		mockInstance.EXPECT().
			NewClient(gomock.Any()).
			Return(mockKafkaClient, nil)

		rgsSchema, err := schrgs.NewSchemaRegistery("something", mockInstance)
		assert.NoError(t, err)

		md, err := rgsSchema.GetLatestSchemaRegistery("something-subject")

		assert.Error(t, err)
		assert.Equal(t, schemaregistry.SchemaMetadata{}, md)
	})
}

func TestSchemaRegistery_CreateAvroSchema(t *testing.T) {
	t.Run("CreateAvroSchema Error Register Schema", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstance := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockClient := mocks.NewMockClient(ctrl)

		mockInstance.EXPECT().
			NewClient(gomock.Any()).
			Return(mockClient, nil)

		mockClient.EXPECT().
			Register(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(0, errors.New("Failed to Register schema"))

		schema, err := schrgs.NewSchemaRegistery("something", mockInstance)
		assert.NoError(t, err)

		id, err := schema.CreateAvroSchema("something", "schema json", false)
		assert.Error(t, err)
		assert.Equal(t, id, 0)
	})

	t.Run("CreateAvroSchema Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstance := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockClient := mocks.NewMockClient(ctrl)

		mockInstance.EXPECT().
			NewClient(gomock.Any()).
			Return(mockClient, nil)

		mockClient.EXPECT().
			Register(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(1, nil)

		schema, err := schrgs.NewSchemaRegistery("something", mockInstance)
		assert.NoError(t, err)

		id, err := schema.CreateAvroSchema("something", "schema json", false)
		assert.NoError(t, err)
		assert.Equal(t, id, 1)
	})
}

func TestGetLatestSchemaId(t *testing.T) {
	t.Run("GetLatestSchemaId Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstance := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockClient := mocks.NewMockClient(ctrl)

		mockInstance.EXPECT().
			NewClient(gomock.Any()).
			Return(mockClient, nil)

		mockClient.EXPECT().
			GetLatestSchemaMetadata(gomock.Any()).
			Return(schemaregistry.SchemaMetadata{}, errors.New("Schema Is Not Found"))

		schema, err := schrgs.NewSchemaRegistery("something", mockInstance)
		assert.NoError(t, err)

		id, err := schema.GetLatestSchemaID("something")
		assert.Error(t, err)
		assert.Equal(t, id, -1)
	})

	t.Run("GetLatestSchemaId", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstance := mocks.NewMockSchemaRegisteryInstance(ctrl)
		mockClient := mocks.NewMockClient(ctrl)

		mockInstance.EXPECT().
			NewClient(gomock.Any()).
			Return(mockClient, nil)

		mockClient.EXPECT().
			GetLatestSchemaMetadata(gomock.Any()).
			Return(schemaregistry.SchemaMetadata{ID: 1}, nil)

		schema, err := schrgs.NewSchemaRegistery("something", mockInstance)
		assert.NoError(t, err)

		id, err := schema.GetLatestSchemaID("something")
		assert.NoError(t, err)
		assert.Equal(t, id, 1)
	})
}
