package schrgs_test

import (
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
	mocks "github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/test/mocks/schrgs"
)

func TestNewSchemaRegisteryClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClient(ctrl)

	client := schrgs.NewSchemaRegisteryClient(mockClient)

	if client.Client() != mockClient {
		t.Error("Expected client to be set correctly")
	}
}

func TestSchemaRegisteryClient_GetLatestSchemaMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClient(ctrl)
	client := schrgs.NewSchemaRegisteryClient(mockClient)

	t.Run("success", func(t *testing.T) {
		subject := "test-subject"
		expectedMetadata := schemaregistry.SchemaMetadata{
			ID:      1,
			Subject: subject,
			Version: 1,
		}

		mockClient.EXPECT().
			GetLatestSchemaMetadata(subject).
			Return(expectedMetadata, nil).
			Times(1)

		result, err := client.GetLatestSchemaMetadata(subject)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result.ID != expectedMetadata.ID {
			t.Errorf("Expected ID %d, got %d", expectedMetadata.ID, result.ID)
		}
	})

	t.Run("error", func(t *testing.T) {
		subject := "test-subject"
		expectedError := errors.New("schema not found")

		mockClient.EXPECT().
			GetLatestSchemaMetadata(subject).
			Return(schemaregistry.SchemaMetadata{}, expectedError).
			Times(1)

		_, err := client.GetLatestSchemaMetadata(subject)

		if err != expectedError {
			t.Errorf("Expected error %v, got %v", expectedError, err)
		}
	})
}

func TestSchemaRegisteryClient_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClient(ctrl)
	client := schrgs.NewSchemaRegisteryClient(mockClient)

	t.Run("success", func(t *testing.T) {
		subject := "test-subject"
		schema := schemaregistry.SchemaInfo{
			Schema:     `{"type": "record"}`,
			SchemaType: "AVRO",
		}
		normalize := true
		expectedID := 123

		mockClient.EXPECT().
			Register(subject, schema, normalize).
			Return(expectedID, nil).
			Times(1)

		id, err := client.Register(subject, schema, normalize)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if id != expectedID {
			t.Errorf("Expected ID %d, got %d", expectedID, id)
		}
	})

	t.Run("error", func(t *testing.T) {
		subject := "test-subject"
		schema := schemaregistry.SchemaInfo{
			Schema:     `invalid schema`,
			SchemaType: "AVRO",
		}
		normalize := false
		expectedError := errors.New("registration failed")

		mockClient.EXPECT().
			Register(subject, schema, normalize).
			Return(0, expectedError).
			Times(1)

		id, err := client.Register(subject, schema, normalize)

		if err != expectedError {
			t.Errorf("Expected error %v, got %v", expectedError, err)
		}
		if id != 0 {
			t.Errorf("Expected ID 0, got %d", id)
		}
	})
}

func TestSchemaRegisteryClient_Client(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClient(ctrl)
	client := schrgs.NewSchemaRegisteryClient(mockClient)

	result := client.Client()

	if result != mockClient {
		t.Error("Expected to return the underlying client")
	}
}
