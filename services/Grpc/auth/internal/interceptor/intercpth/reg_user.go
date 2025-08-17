package intercpth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/validator"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func InterceptRegisterUser(ctx context.Context, sp trace.Span, lg *logs.Logger, req any) error {
	fullMethodName := pbgen.AuthService_RegisterUser_FullMethodName
	code := codes.InvalidArgument
	recorder := NewintercpthErrorRecorder(sp, lg)

	if req == nil {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			"[AuthService] Nil request payload for Register - Expected Request is not nil",
		)
	}

	dataRequest, ok := req.(*pbgen.RegisterUserRequest)
	if !ok {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			fmt.Sprintf("[AuthService] Invalid request type for Register - got: %T - Expected Request have type RegisterRequest Proto", req),
		)
	}

	if !validator.ValidateEmail(dataRequest.Email) {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			fmt.Sprintf("[AuthService] Invalid request type for Register - Email Invalid - %s", dataRequest.Email),
		)
	}

	if !validator.ValidatePhone(dataRequest.PhoneNumber) {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			fmt.Sprintf("[AuthService] Invalid request type for Register - Phone Number Invalid - %s", dataRequest.PhoneNumber),
		)
	}

	if !validator.ValidatePassword(dataRequest.Password) {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			"[AuthService] Invalid request type for Register - Password Invalid - does not meet complexity requirements - Password must be at least 8 characters, include 1 uppercase letter, 1 number, and 1 special character",
		)
	}

	lg.Info(
		ctx,
		fmt.Sprintf("[AuthService] Register request - Email: %s, Phone: %s", dataRequest.Email, dataRequest.PhoneNumber),
		slog.String("full_method", fullMethodName),
		slog.Time("timestamp", time.Now()),
		slog.String("function", "InterceptRegisterUser"),
	)

	sp.AddEvent("validation_completed",
		trace.WithAttributes(
			attribute.String("user.email", dataRequest.GetEmail()),
			attribute.String("user.full_name", dataRequest.GetFullName()),
			attribute.String("user.phone", dataRequest.GetPhoneNumber()),
			attribute.String("layer", "interceptor"),
			attribute.String("validation_status", "success"),
		),
	)
	sp.SetStatus(otelCodes.Ok, "Request validation successful")

	return nil
}
