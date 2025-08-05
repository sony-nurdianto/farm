package unit_test

import (
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/sony-nurdianto/farm/shared_lib/Go/mykafka/pkg"
	"github.com/stretchr/testify/mock"
)

type MockSchemaRegistery struct {
	mock.Mock
}

func (m *MockSchemaRegistery) NewClient(conf *schemaregistry.Config) (schemaregistry.Client, error) {
	args := m.Called(conf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(schemaregistry.Client), args.Error(1)
}

func TestNewSchemaRegistery(t *testing.T) {
	mockSchemaRgs := new(MockSchemaRegistery)
	mockClient := new(MockClient) // << satu-satunya

	mockSchemaRgs.On("NewClient", mock.AnythingOfType("*schemaregistry.Config")).Return(mockClient, nil)

	result, err := pkg.NewSchemaRegistery("something", mockSchemaRgs)
	if err != nil {
		t.Error("Expected NewSchemaRegistery Return Error but got nil")
	}

	t.Run("GetLatestSchemaRegistery", func(t *testing.T) {
		mockClient.On("GetLatestSchemaMetadata", mock.AnythingOfType("string")).
			Return(schemaregistry.SchemaMetadata{}, errors.New("Something Amish"))

		_, err := result.GetLatestSchemaRegistery("something")
		if err == nil {
			t.Error("Expected GetLatestSchemaMetadata Method Return Error")
		}
	})

	t.Run("CreateSchema", func(t *testing.T) {
		mockClient.On("Register",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("schemaregistry.SchemaInfo"),
			false,
		).
			Return(0, errors.New("Something Amish"))

		_, err := result.CreateAvroSchema("something", "somthing", false)
		if err == nil {
			t.Error("Expected GetLatestSchemaMetadata Method Return Error")
		}
	})

	t.Run("GetLatestSchemaID", func(t *testing.T) {
		mockClient.On("GetLatestSchemaMetadata", mock.AnythingOfType("string")).
			Return(schemaregistry.SchemaMetadata{}, errors.New("Something Amish"))

		_, err := result.GetLatestSchemaID("something")
		if err == nil {
			t.Error("Expected GetLatestSchemaMetadata Method Return Error")
		}
	})

	t.Run("Client", func(t *testing.T) {
		c := result.Client()
		if c == nil {
			t.Error("Expected Client Method didn't return nil")
		}
	})
}

func TestNewSchemaRegistery_Error(t *testing.T) {
	mockSchemaRgs := new(MockSchemaRegistery)

	mockSchemaRgs.On("NewClient", mock.AnythingOfType("*schemaregistry.Config")).Return(nil, errors.New("error when creating new shcema"))

	_, err := pkg.NewSchemaRegistery("something", mockSchemaRgs)
	if err == nil {
		t.Error("Expected NewSchemaRegistery Return Error but got nil")
	}
}
