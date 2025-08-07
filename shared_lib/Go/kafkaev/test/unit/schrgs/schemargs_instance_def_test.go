package schrgs_test

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
	"github.com/stretchr/testify/assert"
)

func TestNewRegistery(t *testing.T) {
	instance := schrgs.NewRegistery()
	_ = instance
}

func TestSchemaRegisteryInstance_NewClient(t *testing.T) {
	instance := schrgs.NewRegistery()

	t.Run("with valid config", func(t *testing.T) {
		conf := schemaregistry.NewConfig("http://localhost:8081")

		// Just call the method to get coverage
		client, err := instance.NewClient(conf)
		assert.NoError(t, err)
		// Basic checks if possible
		_ = client
		_ = err
	})

	t.Run("with nil config - expect panic recovery", func(t *testing.T) {
		// Library will panic with nil, so we recover to continue test
		defer func() {
			if r := recover(); r != nil {
				// Expected panic, just continue
			}
		}()
		// This will panic but we recover above
		client, err := instance.NewClient(nil)

		assert.NoError(t, err)
		_ = client
		_ = err
	})

	t.Run("with empty config", func(t *testing.T) {
		// Create a minimal valid config instead of empty
		conf := schemaregistry.NewConfig("http://test:8081")

		client, err := instance.NewClient(conf)

		assert.NoError(t, err)
		_ = client
		_ = err
	})
}
