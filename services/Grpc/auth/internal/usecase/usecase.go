package usecase

import (
	"context"
	"errors"

	"github.com/sony-nurdianto/farm/auth/internal/encryption/passencrypt"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/token"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
)

var (
	ErrorUserIsExist           error = errors.New("User Is Exist Aborting CreateUser")
	ErrorFailedToHasshPassword error = errors.New("Failed To HashPassword")
	ErrorRegisterUser          error = errors.New("Failed To CreateUserAsync")
	ErrorUserIsNotExsist       error = errors.New("User Is Not Exist")
	ErrorPasswordIsInvalid     error = errors.New("Invalid Password Credentials")
)

//go:generate mockgen -package=mocks -destination=../../test/mocks/mock_usecase.go -source=usecase.go
type ServiceUsecase interface {
	UserRegister(ctx context.Context, user *pbgen.RegisterUserRequest) (*pbgen.RegisterUserResponse, error)
	UserSignIn(ctx context.Context, req *pbgen.AuthenticateUserRequest) (*pbgen.AuthenticateUserResponse, error)
	TokenValidate(ctx context.Context, req *pbgen.TokenValidateRequest) (*pbgen.TokenValidateResponse, error)
}

type serviceUsecase struct {
	authRepo    repository.AuthRepo
	passEncrypt passencrypt.PassEncrypt
	tokhen      token.Tokhan
}

func NewServiceUsecase(
	repo repository.AuthRepo,
	pass passencrypt.PassEncrypt,
	tokhen token.Tokhan,
) ServiceUsecase {
	return serviceUsecase{
		authRepo:    repo,
		passEncrypt: pass,
		tokhen:      tokhen,
	}
}
