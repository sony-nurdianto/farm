package usecase

import (
	"context"
	"errors"

	"github.com/sony-nurdianto/farm/auth/internal/encryption/token"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (su serviceUsecase) TokenValidate(
	ctx context.Context,
	req *pbgen.TokenValidateRequest,
) (*pbgen.TokenValidateResponse, error) {
	tracer := otel.Tracer("auth-service")
	_, span := tracer.Start(ctx, "Usecase:TokenValidate")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", "validate_token"),
		attribute.String("layer", "usecase"),
	)

	res := &pbgen.TokenValidateResponse{}

	span.AddEvent("verify_web_token")
	value, err := su.tokhen.VerifyWebToken(req.GetToken())
	if errors.Is(err, token.ErrTokenExperied) {
		res.Valid = false
		res.Msg = "Token Experied"
		span.SetStatus(codes.Error, "User Token Is Experied")
		return res, nil
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed To VerifyWebToken")
		return nil, err
	}

	res.Valid = true
	res.ExpiresAt = timestamppb.New(value.Expiration)
	res.Isuer = &value.Issuer
	res.Subject = &value.Subject
	res.Msg = "Token Is Valid"

	span.AddEvent("verify_token_completed")
	span.SetStatus(codes.Ok, "Token Is Valid And Verify")

	return res, nil
}
