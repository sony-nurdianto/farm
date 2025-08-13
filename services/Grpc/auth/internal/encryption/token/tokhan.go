package token

import (
	"errors"
	"os"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

var (
	ErrSecretNotSet  error = errors.New("SECRET is not set")
	ErrDecryptFailed error = errors.New("invalid token failed to decrypt")
	ErrTokenExperied error = errors.New("invalid token token is experied")
)

//go:generate mockgen -source=tokhan.go -destination=../../../test/mocks/mock_tokhan.go -package=mocks
type Tokhan interface {
	CreateWebToken(subject string) (string, error)
	VerifyWebToken(token string) (paseto.JSONToken, error)
}

type tokhan struct {
	handler PassetoToken
}

func NewTokhan(handler PassetoToken) Tokhan {
	return tokhan{handler}
}

func secret() []byte {
	secretKey := make([]byte, chacha20poly1305.KeySize)
	secret := os.Getenv("TOKSECRET")

	copy(secretKey, secret)
	return secretKey
}

func (tg tokhan) CreateWebToken(subject string) (string, error) {
	jsonToken := paseto.JSONToken{
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(1 * time.Hour),
		Subject:    subject,
		Issuer:     "auth",
	}

	token, err := tg.handler.Encrypt(secret(), jsonToken, nil)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (tg tokhan) VerifyWebToken(token string) (paseto.JSONToken, error) {
	var newJsonToken paseto.JSONToken

	if err := tg.handler.Decrypt(token, secret(), &newJsonToken, nil); err != nil {
		return newJsonToken, ErrDecryptFailed
	}

	if newJsonToken.Expiration.Before(time.Now()) {
		return newJsonToken, ErrTokenExperied
	}

	return newJsonToken, nil
}
