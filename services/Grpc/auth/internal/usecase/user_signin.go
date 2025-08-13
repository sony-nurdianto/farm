package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
