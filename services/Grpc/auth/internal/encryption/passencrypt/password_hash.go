package passencrypt

import (
	"errors"
	"fmt"
	"io"

	"github.com/sony-nurdianto/farm/auth/internal/encryption/codec"
	"golang.org/x/crypto/argon2"
)

//go:generate mockgen -source=password_hash.go -destination=../../../test/mocks/mock_passencrypt.go -package=mocks
const (
	time    = 1
	memory  = 64 * 1024
	threads = 4
	keyLen  = 32
)

var ErrorGenSaltReadRand = errors.New("failed to read rand")

type PassEncrypt interface {
	HashPassword(password string) (passwordHash string, _ error)
	VerifyPassword(password, passwordHash string) (verify bool, _ error)
}

type passEncrypt struct {
	randRead       io.Reader
	base64Encoding codec.Base64Encoder
}

func NewPassEncrypt(randRead io.Reader, base64Encoding codec.Base64Encoder) passEncrypt {
	return passEncrypt{randRead, base64Encoding}
}

func (pe passEncrypt) HashPassword(
	password string,
) (passwordHash string, _ error) {
	salt := make([]byte, 32)

	_, err := pe.randRead.Read(salt)
	if err != nil {
		return passwordHash, err
	}

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	saltStr := pe.base64Encoding.EncodeToString(salt)
	hashStr := pe.base64Encoding.EncodeToString(hash)

	passwordHash = fmt.Sprintf("%s%s", hashStr, saltStr)

	return passwordHash, nil
}

func (pe passEncrypt) VerifyPassword(
	password, passwordHash string,
) (verify bool, _ error) {
	salt, err := pe.base64Encoding.DecodeString(passwordHash[44:])
	if err != nil {
		return verify, err
	}

	expectedHash, err := pe.base64Encoding.DecodeString(passwordHash[:44])
	if err != nil {
		return verify, err
	}

	newHash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	verify = string(newHash) == string(expectedHash)

	return verify, nil
}
