package unit_test

import (
	"crypto/rand"
	"errors"
	"strings"
	"testing"

	"github.com/sony-nurdianto/farm/auth/internal/encryption/codec"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/passencrypt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBase64Encoder struct {
	mock.Mock
}

func (cd *MockBase64Encoder) EncodeToString(src []byte) string {
	args := cd.Mock.Called(src)
	return args.Get(0).(string)
}

func (cd *MockBase64Encoder) DecodeString(s string) ([]byte, error) {
	args := cd.Mock.Called(s)
	return args.Get(0).([]byte), args.Error(1)
}

type MockReader struct {
	mock.Mock
}

func (m *MockReader) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func TestPasswordHash_Error(t *testing.T) {
	mockCodec := new(MockBase64Encoder)
	mockReader := new(MockReader)

	t.Run("Error Rand Read", func(t *testing.T) {
		errMsg := errors.New("Failed to Read")

		mockReader.On("Read", mock.AnythingOfType("[]uint8")).Return(0, errMsg)
		mockCodec.On("EncodeToString", mock.AnythingOfType("[]uint8")).Return("", nil)

		pe := passencrypt.NewPassEncrypt(mockReader, mockCodec)

		_, err := pe.HashPassword("password")
		if err == nil {
			assert.Error(t, errMsg)
		}
	})
}

func TestPasswordVerify_Error(t *testing.T) {
	mockReader := new(MockReader)

	t.Run("Error DecodeString Salt", func(t *testing.T) {
		mockCodec := new(MockBase64Encoder)

		errMsg := errors.New("Failed to Read")

		mockReader.On("Read", mock.AnythingOfType("[]uint8")).Return(32, nil)
		mockCodec.On("DecodeString", mock.AnythingOfType("string")).Return([]byte{}, errMsg).Once()

		pe := passencrypt.NewPassEncrypt(mockReader, mockCodec)

		_, err := pe.VerifyPassword("password", strings.Repeat("h", 88))
		if err == nil {
			assert.Error(t, errMsg)
		}

		mockCodec.AssertNumberOfCalls(t, "DecodeString", 1)
		mockCodec.AssertExpectations(t)
	})

	t.Run("Error DecodeString Hash", func(t *testing.T) {
		mockCodec := new(MockBase64Encoder)
		errMsg := errors.New("Failed to Read")

		mockReader.On("Read", mock.AnythingOfType("[]uint8")).Return(32, nil)
		mockCodec.On("DecodeString", mock.AnythingOfType("string")).Return([]byte{}, nil).Once()
		mockCodec.On("DecodeString", mock.AnythingOfType("string")).Return([]byte{}, errMsg).Once()

		pe := passencrypt.NewPassEncrypt(mockReader, mockCodec)

		_, err := pe.VerifyPassword("password", strings.Repeat("h", 88))
		if err == nil {
			assert.Error(t, errMsg)
		}

		mockCodec.AssertNumberOfCalls(t, "DecodeString", 2)
		mockCodec.AssertExpectations(t)
	})
}

func TestPassEncrypt(t *testing.T) {
	codec := codec.NewBase64Encoder()

	password := "secret"

	pe := passencrypt.NewPassEncrypt(rand.Reader, codec)

	passwordHash, err := pe.HashPassword(password)
	if err != nil {
		t.Errorf("Expected HashPassword Method didn't return err but got: %s", err)
	}

	verify, err := pe.VerifyPassword(password, passwordHash)
	if err != nil {
		t.Errorf("Expected VerifyPassword Method didn't return err but got: %s", err)
	}

	if !verify {
		t.Errorf("Expected verify is true but got %t", verify)
	}
}
