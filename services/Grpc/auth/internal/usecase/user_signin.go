package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (su serviceUsecase) UserSignIn(ctx context.Context, req *pbgen.AuthenticateUserRequest) (*pbgen.AuthenticateUserResponse, error) {
	tracer := otel.Tracer("auth-service")
	uctx, span := tracer.Start(ctx, "Usecase:UserSignIn")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", "user_signin"),
		attribute.String("layer", "usecase"),
	)

	span.AddEvent("get_user_by_email")
	user, err := su.authRepo.GetUserByEmail(uctx, req.GetEmail())
	if errors.Is(err, sql.ErrNoRows) {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("User with email %s not found", req.GetEmail()))
		return nil, fmt.Errorf("%w: user with email %s is not exist", ErrorUserIsNotExsist, req.GetEmail())
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed Get User By Email")
		return nil, err
	}

	span.AddEvent("verify_user_password")
	isPass, err := su.passEncrypt.VerifyPassword(req.Password, user.Password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed Verify User Password")
		return nil, err
	}

	if !isPass {
		span.RecordError(errors.New("Unauthorized Password is Invalid"))
		span.SetStatus(codes.Error, "User Password Is Invalid")
		return nil, ErrorPasswordIsInvalid
	}

	span.AddEvent("create_user_token")
	createToken, err := su.tokhen.CreateWebToken(user.Id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed Create User Token")
		return nil, err
	}

	response := &pbgen.AuthenticateUserResponse{
		Token:     createToken,
		Status:    "Success",
		Msg:       "User Authenticated Success Login. Welcome !",
		IssuedAt:  timestamppb.Now(),
		ExpiresAt: timestamppb.New(time.Now().Add(1 * time.Hour)),
	}

	span.AddEvent("user_signin_completed")
	span.SetStatus(codes.Ok, "User SignIn successfully")

	return response, nil
}
