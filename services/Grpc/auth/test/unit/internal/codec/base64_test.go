package unit_test

import (
	"testing"

	"github.com/sony-nurdianto/farm/auth/internal/encryption/codec"
	"github.com/stretchr/testify/assert"
)

func TestBase64Codec(t *testing.T) {
	codec := codec.NewBase64Encoder()
	en := codec.EncodeToString([]byte("Something"))
	assert.NotEmpty(t, en)

	de, err := codec.DecodeString(en)
	assert.NoError(t, err)
	assert.Equal(t, string(de), "Something")
}
