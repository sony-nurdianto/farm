package usecase

import (
	"crypto/rand"
	"database/sql"
	"errors"

	"github.com/sony-nurdianto/farm/auth/internal/encryption/codec"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/passencrypt"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
)

var ErrorUserIsExist error = errors.New("User Is Exist Aborting CreateUser")

type ServiceUsecase struct {
	RepoPG *repository.RepoPostgres
}

func checkUser(rp *repository.RepoPostgres, email string) (bool, error) {
	_, err := rp.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (su ServiceUsecase) UserRegister(user *pbgen.RegisterRequest) (*pbgen.RegisterResponse, error) {
	userExsist, err := checkUser(su.RepoPG, user.GetEmail())
	if err != nil {
		return nil, err
	}

	if userExsist {
		return nil, ErrorUserIsExist
	}

	pe := passencrypt.NewPassEncrypt(
		rand.Reader,
		codec.NewBase64Encoder(),
	)

	passwordHash, err := pe.HashPassword(user.GetPassword())
	if err != nil {
		return nil, err
	}

	_, err = su.RepoPG.CreateUser(user.GetEmail(), passwordHash)
	if err != nil {
		return nil, err
	}

	out := &pbgen.RegisterResponse{
		Msg:    "Success Create User",
		Status: "Success",
	}

	return out, nil
}
