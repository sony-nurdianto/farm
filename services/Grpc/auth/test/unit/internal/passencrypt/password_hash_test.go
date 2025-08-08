package unit_test

import (
	"crypto/rand"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/codec"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/passencrypt"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/stretchr/testify/assert"
)

type errReader struct{}

func (e errReader) Read(p []byte) (int, error) {
	return 0, errors.New("mocked read error")
}

func TestPasswordHash(t *testing.T) {
	t.Run("Hash And Verify Password", func(t *testing.T) {
		codec := codec.NewBase64Encoder()
		password := "Password"
		passEncrypt := passencrypt.NewPassEncrypt(rand.Reader, codec)

		res, err := passEncrypt.HashPassword("Password")
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		ok, err := passEncrypt.VerifyPassword(password, res)
		assert.NoError(t, err)
		assert.True(t, ok)

		wrong, err := passEncrypt.VerifyPassword("My Password Is Wrong", res)
		assert.NoError(t, err)
		assert.False(t, wrong)
	})

	t.Run("Hash Password Error Reader Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCodec := mocks.NewMockBase64Encoder(ctrl)

		passEncrypt := passencrypt.NewPassEncrypt(errReader{}, mockCodec)

		hp, err := passEncrypt.HashPassword("Password")
		assert.Error(t, err)
		assert.Empty(t, hp)
	})

	t.Run("VerifyPassword Error DecodeString Salt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCodec := mocks.NewMockBase64Encoder(ctrl)
		mockCodec.EXPECT().
			DecodeString(gomock.Any()).
			Return([]byte{}, errors.New("Failed to Decode Salt String")).
			Times(1)

		passencrypt := passencrypt.NewPassEncrypt(rand.Reader, mockCodec)

		res, err := passencrypt.VerifyPassword("Password", strings.Repeat("x", 88))
		assert.Error(t, err)
		assert.False(t, res)
	})

	t.Run("VerifyPassword Error DecodeString Hash", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCodec := mocks.NewMockBase64Encoder(ctrl)
		mockCodec.EXPECT().
			DecodeString(gomock.Any()).
			Return([]byte{}, nil).
			Times(1)

		mockCodec.EXPECT().
			DecodeString(gomock.Any()).
			Return([]byte{}, errors.New("Failed to Decode Hash String")).
			Times(1)

		passencrypt := passencrypt.NewPassEncrypt(rand.Reader, mockCodec)

		res, err := passencrypt.VerifyPassword("Password", strings.Repeat("x", 88))
		assert.Error(t, err)
		assert.EqualError(t, err, "Failed to Decode Hash String")
		assert.False(t, res)
	})
}
