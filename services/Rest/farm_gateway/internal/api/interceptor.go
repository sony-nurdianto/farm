package api

import (
	"context"
	"log"

	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"

	"google.golang.org/grpc/status"
)

type unaryClientInterceptor struct {
	traceProvider *trace.TraceProvider
	logger        *logs.Logger
}

func NewUnaryClientInterceptor(tp *trace.TraceProvider) unaryClientInterceptor {
	return unaryClientInterceptor{
		traceProvider: tp,
		logger:        logs.NewLogger(),
	}
}

func (uci unaryClientInterceptor) UnaryAuthClientIntercept(
	ctx context.Context,
	method string,
	req, res any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	tracer := uci.traceProvider.Tracer("auth-service-client")

	ctx, span := tracer.Start(ctx, method)
	defer span.End()

	// Tambahkan attribute awal
	span.SetAttributes(
		attribute.String("rpc.system", "grpc"),
		attribute.String("rpc.method", method),
	)

	err := invoker(ctx, method, req, res, cc, opts...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.SetAttributes(attribute.String("rpc.status_code", status.Code(err).String()))

	log.Println("gRPC method:", method, "response:", res, "error:", err)
	return err
}
