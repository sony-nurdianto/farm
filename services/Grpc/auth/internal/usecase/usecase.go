package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/passencrypt"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/token"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	UserRegister(user *pbgen.RegisterUserRequest) (*pbgen.RegisterUserResponse, error)
	UserSignIn(req *pbgen.AuthenticateUserRequest) (*pbgen.AuthenticateUserResponse, error)
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

func (su serviceUsecase) UserRegister(user *pbgen.RegisterUserRequest) (*pbgen.RegisterUserResponse, error) {
	userExsist, err := checkUser(su.authRepo, user.GetEmail())
	if err != nil {
		return nil, err
	}

	if userExsist {
		return nil, fmt.Errorf("%w: %s", ErrorUserIsExist, err)
	}

	passwordHash, err := su.passEncrypt.HashPassword(user.GetPassword())
	if err != nil {
		return nil, ErrorFailedToHasshPassword
	}

	userId := uuid.NewString()

	err = su.authRepo.CreateUserAsync(userId, user.GetEmail(), user.GetFullName(), user.GetPhoneNumber(), passwordHash)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrorRegisterUser, err)
	}

	out := &pbgen.RegisterUserResponse{
		Msg:    "Success Create User",
		Status: "Success",
	}

	return out, nil
}

func (su serviceUsecase) UserSignIn(req *pbgen.AuthenticateUserRequest) (*pbgen.AuthenticateUserResponse, error) {
	user, err := su.authRepo.GetUserByEmail(req.GetEmail())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: user with email %s is not exist", ErrorUserIsNotExsist, req.GetEmail())
	}

	if err != nil {
		return nil, err
	}

	isPass, err := su.passEncrypt.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return nil, err
	}

	if !isPass {
		return nil, ErrorPasswordIsInvalid
	}

	createToken, err := su.tokhen.CreateWebToken(user.Id)
	if err != nil {
		return nil, err
	}

	response := &pbgen.AuthenticateUserResponse{
		Token:     createToken,
		Status:    "Success",
		Msg:       "User Authenticated Success Login. Welcome !",
		IssuedAt:  timestamppb.Now(),
		ExpiresAt: timestamppb.New(time.Now().Add(1 * time.Hour)),
	}

	return response, nil
}
