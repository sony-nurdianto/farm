package codec

import "encoding/base64"

//go:generate mockgen -source=base64.go -destination=../../../test/mocks/mock_codec.go -package=mocks
type Base64Encoder interface {
	EncodeToString(src []byte) string
	DecodeString(string) ([]byte, error)
}

type base64encoder struct{}

func NewBase64Encoder() base64encoder {
	return base64encoder{}
}

func (cd base64encoder) EncodeToString(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func (cd base64encoder) DecodeString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
