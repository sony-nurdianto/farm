package interceptor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	"google.golang.org/grpc"
)

func AuthServiceUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	logger := logs.NewLogger()

	logger.Info(
		ctx,
		fmt.Sprintf("[FarmerService] Incoming request - Method: %s", info.FullMethod),
		slog.String("full_method", info.FullMethod),
	)

	resp, err = handler(ctx, req)

	return resp, err
}
