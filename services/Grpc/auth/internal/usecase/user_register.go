package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func checkUser(ctx context.Context, rp repository.AuthRepo, email string) (bool, error) {
	_, err := rp.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (su serviceUsecase) UserRegister(ctx context.Context, user *pbgen.RegisterUserRequest) (*pbgen.RegisterUserResponse, error) {
	tracer := otel.Tracer("auth-service")
	uctx, span := tracer.Start(ctx, "Usecase:UserRegister")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", "user_registration"),
		attribute.String("layer", "usecase"),
	)

	span.AddEvent("checking_user_existence")
	userExists, err := checkUser(uctx, su.authRepo, user.GetEmail())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to check user existence")
		return nil, err
	}

	if userExists {
		err := fmt.Errorf("%w: user already exists", ErrorUserIsExist)
		span.RecordError(err)
		span.SetStatus(codes.Error, "User already exists")
		return nil, err
	}

	span.AddEvent("hashing_password")
	passwordHash, err := su.passEncrypt.HashPassword(user.GetPassword())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to hash password")
		return nil, ErrorFailedToHasshPassword
	}

	span.AddEvent("creating_user")
	userId := uuid.NewString()
	span.SetAttributes(attribute.String("user.id", userId))

	err = su.authRepo.CreateUserAsync(
		uctx,
		userId,
		user.GetEmail(),
		user.GetFullName(),
		user.GetPhoneNumber(),
		passwordHash,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create user")
		return nil, fmt.Errorf("%w: %s", ErrorRegisterUser, err)
	}

	span.AddEvent("user_registration_completed")
	span.SetStatus(codes.Ok, "User registered successfully")

	out := &pbgen.RegisterUserResponse{
		Msg:    "Success Create User",
		Status: "Success",
	}
	return out, nil
}
