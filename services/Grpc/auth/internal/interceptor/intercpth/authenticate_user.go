package intercpth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/recorderr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
)

func InterceptAuthenticateUser(ctx context.Context, sp trace.Span, lg *logs.Logger, req any) error {
	fullMethodName := pbgen.AuthService_AuthenticateUser_FullMethodName
	code := codes.InvalidArgument
	recorder := recorderr.NewErrorRecorder(sp, lg)
	if req == nil {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			"[AuthService] Nil request payload for AuthenticateUser - Expected Request is not nil",
		)
	}

	dataRequest, ok := req.(*pbgen.AuthenticateUserRequest)
	if !ok {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			fmt.Sprintf("[AuthService] Invalid request type for AuthenticateUser - got: %T - Expected Request have type AuthenticateUserRequest Proto", req),
		)
	}

	if len(dataRequest.GetEmail()) == 0 {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			"[AuthService] Invalid request type for AuthenticateUser - Email is empty - does not requirements - Email must not be empty",
		)
	}

	if len(dataRequest.GetPassword()) == 0 {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			"[AuthService] Invalid request type for AuthenticateUser - Password is empty - does not requirements - Password must not be empty",
		)
	}

	lg.Info(
		ctx,
		fmt.Sprintf("[AuthService] AuthenticateUser request - Email: %s", dataRequest.Email),
		slog.String("full_method", fullMethodName),
		slog.Time("timestamp", time.Now()),
		slog.String("function", "InterceptAuthenticateUser"),
	)

	sp.AddEvent("validation_completed",
		trace.WithAttributes(
			attribute.String("email", dataRequest.Email),
			attribute.String("validation_status", "success"),
		),
	)
	sp.SetStatus(otelCodes.Ok, "Request validation successful")

	return nil
}
