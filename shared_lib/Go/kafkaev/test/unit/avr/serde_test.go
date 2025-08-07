package unit_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	mocks "github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/test/mocks/avr"
	mockClient "github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/test/mocks/schrgs"
	"github.com/stretchr/testify/assert"
)

func TestNewAvroGenericSerde(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("NewAvroGenericSerde", func(t *testing.T) {
		mocksSchemaRgsClient := mockClient.NewMockClient(ctrl)

		mockAvrSerializer := mocks.NewMockAvrSerializer(ctrl)
		mockAvrDeserilaizer := mocks.NewMockAvrDeserializer(ctrl)

		mockAvrInstance := mocks.NewMockAvrSerdeInstance(ctrl)
		mockAvrInstance.EXPECT().
			NewGenericSerializer(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockAvrSerializer, nil)

		mockAvrInstance.EXPECT().
			NewGenericDeserializer(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockAvrDeserilaizer, nil)

		_, err := avr.NewAvroGenericSerde(mocksSchemaRgsClient, mockAvrInstance)
		assert.NoError(t, err)
	})

	t.Run("NewAvroGenericSerde Error Create NewGenericSerializer", func(t *testing.T) {
		mocksSchemaRgsClient := mockClient.NewMockClient(ctrl)

		mockAvrSerializer := mocks.NewMockAvrSerializer(ctrl)

		mockAvrInstance := mocks.NewMockAvrSerdeInstance(ctrl)
		mockAvrInstance.EXPECT().
			NewGenericSerializer(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockAvrSerializer, errors.New("Failed To Create New GenericSerializer"))

		_, err := avr.NewAvroGenericSerde(mocksSchemaRgsClient, mockAvrInstance)
		assert.Error(t, err)
	})

	t.Run("NewAvroGenericSerde Error Create NewGenericDeserializer", func(t *testing.T) {
		mocksSchemaRgsClient := mockClient.NewMockClient(ctrl)

		mockAvrSerializer := mocks.NewMockAvrSerializer(ctrl)
		mockAvrDeserilaizer := mocks.NewMockAvrDeserializer(ctrl)

		mockAvrInstance := mocks.NewMockAvrSerdeInstance(ctrl)
		mockAvrInstance.EXPECT().
			NewGenericSerializer(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockAvrSerializer, nil)

		mockAvrInstance.EXPECT().
			NewGenericDeserializer(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockAvrDeserilaizer, errors.New("Failed To Create New GenericDeserializer"))

		_, err := avr.NewAvroGenericSerde(mocksSchemaRgsClient, mockAvrInstance)
		assert.Error(t, err)
	})
}
