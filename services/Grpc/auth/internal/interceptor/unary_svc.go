package interceptor

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/sony-nurdianto/farm/auth/internal/interceptor/intercpth"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/trace"
	"google.golang.org/grpc"
)

func AuthServiceUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	start := time.Now()
	ts := trace.NewTraceSpan()
	span := ts.SpanFromContext(ctx)
	logger := logs.NewLogger()

	logger.Info(
		ctx,
		fmt.Sprintf("[AuthService] Incoming request - Method: %s", info.FullMethod),
		slog.String("full_method", info.FullMethod),
	)

	switch info.FullMethod {
	case pbgen.AuthService_RegisterUser_FullMethodName:
		if err := intercpth.InterceptRegisterUser(ctx, span, logger, req); err != nil {
			return nil, err
		}

	case pbgen.AuthService_AuthenticateUser_FullMethodName:
		if err := intercpth.InterceptAuthenticateUser(ctx, span, logger, req); err != nil {
			return nil, err
		}
	case pbgen.AuthService_TokenValidate_FullMethodName:
		if err := intercpth.InterceptTokenValidate(ctx, span, logger, req); err != nil {
			return nil, err
		}
	}

	resp, err = handler(ctx, req)

	duration := time.Since(start)
	if err != nil {
		// errfac := fmt.Errorf("[AuthService] Method %s failed - Error: %v - Duration: %v", info.FullMethod, err, duration)
	} else {
		log.Printf("[AuthService] Method %s succeeded - Duration: %v", info.FullMethod, duration)
	}

	return resp, err
}
