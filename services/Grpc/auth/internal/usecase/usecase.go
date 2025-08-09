package usecase

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/passencrypt"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
)

var ErrorUserIsExist error = errors.New("User Is Exist Aborting CreateUser")

type ServiceUsecase struct {
	authRepo    repository.AuthRepo
	passEncrypt passencrypt.PassEncrypt
}

func NewServiceUsecase(
	repo repository.AuthRepo,
	pass passencrypt.PassEncrypt,
) ServiceUsecase {
	return ServiceUsecase{
		authRepo:    repo,
		passEncrypt: pass,
	}
}

func checkUser(rp repository.AuthRepo, email string) (bool, error) {
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
	userExsist, err := checkUser(su.authRepo, user.GetEmail())
	if err != nil {
		return nil, err
	}

	if userExsist {
		return nil, ErrorUserIsExist
	}

	passwordHash, err := su.passEncrypt.HashPassword(user.GetPassword())
	if err != nil {
		return nil, err
	}

	userId := uuid.NewString()

	err = su.authRepo.CreateUserAsync(userId, user.GetEmail(), user.GetFullName(), user.GetPhoneNumber(), passwordHash)
	if err != nil {
		return nil, err
	}

	out := &pbgen.RegisterResponse{
		Msg:    "Success Create User",
		Status: "Success",
	}

	return out, nil
}
