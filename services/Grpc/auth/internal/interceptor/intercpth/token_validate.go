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

func InterceptTokenValidate(ctx context.Context, sp trace.Span, lg *logs.Logger, req any) error {
	fullMethodName := pbgen.AuthService_AuthenticateUser_FullMethodName
	code := codes.InvalidArgument
	recorder := recorderr.NewErrorRecorder(sp, lg)
	if req == nil {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			"[AuthService] Nil request payload for TokenValidate - Expected Request is not nil",
		)
	}

	dataRequest, ok := req.(*pbgen.TokenValidateRequest)
	if !ok {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			fmt.Sprintf("[AuthService] Invalid request type for TokenValidate - got: %T - Expected Request have type TokenValidateRequest Proto", req),
		)
	}

	if len(dataRequest.Token) == 0 {
		return recorder.Record(
			ctx,
			code,
			fullMethodName,
			"[AuthService] Invalid request type for TokenValidate - Token is empty - does not meet requirements",
		)
	}

	lg.Info(
		ctx,
		"[AuthService] Token Validate request",
		slog.String("full_method", fullMethodName),
		slog.Time("timestamp", time.Now()),
		slog.String("function", "InterceptTokenValidate"),
	)

	sp.AddEvent("validation_completed",
		trace.WithAttributes(
			attribute.String("validation_status", "success"),
		),
	)
	sp.SetStatus(otelCodes.Ok, "Request validation successful")

	return nil
}
