package usecase

import (
	"context"
	"errors"

	"github.com/sony-nurdianto/farm/auth/internal/encryption/token"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (su serviceUsecase) TokenValidate(
	ctx context.Context,
	req *pbgen.TokenValidateRequest,
) (*pbgen.TokenValidateResponse, error) {
	res := &pbgen.TokenValidateResponse{}

	value, err := su.tokhen.VerifyWebToken(req.GetToken())
	if errors.Is(err, token.ErrTokenExperied) {
		res.Valid = false
		res.Msg = "Token Experied"
		return res, nil
	}

	if err != nil {
		return nil, err
	}

	res.Valid = true
	res.ExpiresAt = timestamppb.New(value.Expiration)
	res.Isuer = &value.Issuer
	res.Subject = &value.Subject
	res.Msg = "Token Is Valid"

	return res, nil
}
