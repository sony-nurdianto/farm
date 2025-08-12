package token

import (
	"github.com/o1egl/paseto"
)

//go:generate mockgen -source=passeto_def.go -destination=../../../test/mocks/mock_passeto_def.go -package=mocks
type PassetoToken interface {
	Encrypt(secretKey []byte, payload any, footer any) (string, error)
	Decrypt(token string, key []byte, payload any, footer any) error
}

type passetoToken struct {
	passetoV2 *paseto.V2
}

func NewPassetoToken() PassetoToken {
	return passetoToken{passetoV2: paseto.NewV2()}
}

func (pt passetoToken) Encrypt(secretKey []byte, payload any, footer any) (string, error) {
	return pt.passetoV2.Encrypt(secretKey, payload, footer)
}

func (pt passetoToken) Decrypt(token string, key []byte, payload any, footer any) error {
	return pt.passetoV2.Decrypt(token, key, payload, footer)
}
