package service

import (
	"context"
	"errors"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/recorderr"
	"github.com/sony-nurdianto/farm/auth/internal/usecase"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"

	otelCodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceServer struct {
	pbgen.UnimplementedAuthServiceServer
	serviceUsecase usecase.ServiceUsecase
}

func NewAuthServiceServer(uc usecase.ServiceUsecase) *AuthServiceServer {
	return &AuthServiceServer{serviceUsecase: uc}
}

func handleRegisterError(ctx context.Context, err error, errRecorder recorderr.ErrorRecorder) error {
	fullMethodName := pbgen.AuthService_RegisterUser_FullMethodName
	switch {
	case err == usecase.ErrorUserIsExist:
		return errRecorder.Record(ctx, codes.AlreadyExists, fullMethodName, err.Error())
	case errors.Is(err, usecase.ErrorFailedToHasshPassword):
		return errRecorder.Record(ctx, codes.Internal, fullMethodName, err.Error())
	case errors.Is(err, usecase.ErrorRegisterUser):
		return errRecorder.Record(ctx, codes.Internal, fullMethodName, err.Error())
	default:
		return errRecorder.Record(ctx, codes.Internal, fullMethodName, err.Error())
	}
}

func (ass *AuthServiceServer) RegisterUser(
	ctx context.Context,
	in *pbgen.RegisterUserRequest,
) (*pbgen.RegisterUserResponse, error) {
	tracer := otel.Tracer("auth-service")
	hctx, span := tracer.Start(ctx, "ServiceHandler:RegisterUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", "user_registration"),
		attribute.String("layer", "handler"),
	)

	errRecorder := recorderr.NewErrorRecorder(span, logs.NewLogger())

	res, err := ass.serviceUsecase.UserRegister(hctx, in)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "User Registration Failed")
		return nil, handleRegisterError(ctx, err, errRecorder)
	}

	span.SetStatus(otelCodes.Ok, "User registration completed successfully")
	return res, nil
}

func handleAutUserErr(err error) error {
	switch {
	case errors.Is(err, usecase.ErrorUserIsNotExsist):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, usecase.ErrorPasswordIsInvalid):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (ass *AuthServiceServer) AuthenticateUser(
	ctx context.Context,
	in *pbgen.AuthenticateUserRequest,
) (*pbgen.AuthenticateUserResponse, error) {
	res, err := ass.serviceUsecase.UserSignIn(ctx, in)
	if err != nil {
		return nil, handleAutUserErr(err)
	}

	return res, nil
}

func (ass *AuthServiceServer) TokenValidate(
	ctx context.Context,
	in *pbgen.TokenValidateRequest,
) (*pbgen.TokenValidateResponse, error) {
	res, err := ass.serviceUsecase.TokenValidate(ctx, in)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return res, nil
}
