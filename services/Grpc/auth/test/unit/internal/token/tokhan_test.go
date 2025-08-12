package unit_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/o1egl/paseto"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/token"
	"github.com/sony-nurdianto/farm/auth/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateWebToken(t *testing.T) {
	t.Run("CreateWebToken Error Encrypt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mocksPasTok := mocks.NewMockPassetoToken(ctrl)

		mocksPasTok.EXPECT().
			Encrypt(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).Return("", errors.New("Failed To Encrypt Token"))

		tokhan := token.NewTokhan(mocksPasTok)
		_, err := tokhan.CreateWebToken("something")
		assert.Error(t, err)
	})

	t.Run("CreateWebToken Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mocksPasTok := mocks.NewMockPassetoToken(ctrl)

		mocksPasTok.EXPECT().
			Encrypt(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).Return("", nil)

		tokhan := token.NewTokhan(mocksPasTok)
		_, err := tokhan.CreateWebToken("something")
		assert.NoError(t, err)
	})
}

func TestVerifyWebToken(t *testing.T) {
	t.Run("VerifyWebToken Error Decrypt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mocksPasTok := mocks.NewMockPassetoToken(ctrl)

		mocksPasTok.EXPECT().
			Decrypt(
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).Return(errors.New("Failed To Encrypt Token"))

		tokhan := token.NewTokhan(mocksPasTok)
		_, err := tokhan.VerifyWebToken("something")
		assert.Error(t, err)
	})

	t.Run("VerifyWebToken Error Token Expired", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expiredToken := paseto.JSONToken{
			Expiration: time.Now().Add(-1 * time.Hour),
			Subject:    "user1",
		}

		mocksPasTok := mocks.NewMockPassetoToken(ctrl)

		mocksPasTok.EXPECT().
			Decrypt(gomock.Any(), gomock.Any(), gomock.Any(), nil).
			DoAndReturn(func(t string, key []byte, payload any, footer any) error {
				p := payload.(*paseto.JSONToken)
				*p = expiredToken
				return nil
			})

		tokhan := token.NewTokhan(mocksPasTok)
		_, err := tokhan.VerifyWebToken("something")
		assert.Error(t, err)
		assert.ErrorIs(t, err, token.ErrTokenExperied)
	})

	t.Run("VerifyWebToken Error Token Expired", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tok := paseto.JSONToken{
			Expiration: time.Now().Add(1 * time.Hour),
			Subject:    "user1",
		}

		mocksPasTok := mocks.NewMockPassetoToken(ctrl)

		mocksPasTok.EXPECT().
			Decrypt(gomock.Any(), gomock.Any(), gomock.Any(), nil).
			DoAndReturn(func(t string, key []byte, payload any, footer any) error {
				p := payload.(*paseto.JSONToken)
				*p = tok
				return nil
			})

		tokhan := token.NewTokhan(mocksPasTok)
		tokDec, err := tokhan.VerifyWebToken("something")
		assert.NoError(t, err)
		assert.Equal(t, tokDec.Subject, tok.Subject)
	})
}
