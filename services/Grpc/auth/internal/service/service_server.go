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
		return nil, handleRegisterError(hctx, err, errRecorder)
	}

	span.SetStatus(otelCodes.Ok, "User registration completed successfully")
	return res, nil
}

func handleAutUserErr(ctx context.Context, err error, errRecorder recorderr.ErrorRecorder) error {
	fullMethodName := pbgen.AuthService_AuthenticateUser_FullMethodName
	switch {
	case errors.Is(err, usecase.ErrorUserIsNotExsist):
		return errRecorder.Record(ctx, codes.NotFound, fullMethodName, err.Error())
	case errors.Is(err, usecase.ErrorPasswordIsInvalid):
		return errRecorder.Record(ctx, codes.InvalidArgument, fullMethodName, err.Error())
	default:
		return errRecorder.Record(ctx, codes.Internal, fullMethodName, err.Error())
	}
}

func (ass *AuthServiceServer) AuthenticateUser(
	ctx context.Context,
	in *pbgen.AuthenticateUserRequest,
) (*pbgen.AuthenticateUserResponse, error) {
	tracer := otel.Tracer("auth-service")
	hctx, span := tracer.Start(ctx, "ServiceHandler:RegisterUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", "authenticate_user"),
		attribute.String("layer", "handler"),
	)

	errRecorder := recorderr.NewErrorRecorder(span, logs.NewLogger())
	res, err := ass.serviceUsecase.UserSignIn(hctx, in)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "User SignIn Failed")
		return nil, handleAutUserErr(hctx, err, errRecorder)
	}

	return res, nil
}

func (ass *AuthServiceServer) TokenValidate(
	ctx context.Context,
	in *pbgen.TokenValidateRequest,
) (*pbgen.TokenValidateResponse, error) {
	tracer := otel.Tracer("auth-service")
	hctx, span := tracer.Start(ctx, "ServiceHandler:RegisterUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", "authenticate_user"),
		attribute.String("layer", "handler"),
	)

	errRecorder := recorderr.NewErrorRecorder(span, logs.NewLogger())
	fullMethodName := pbgen.AuthService_TokenValidate_FullMethodName
	res, err := ass.serviceUsecase.TokenValidate(hctx, in)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "Failed Validate Token")
		return nil, errRecorder.Record(hctx, codes.InvalidArgument, fullMethodName, err.Error())
	}

	return res, nil
}
