package trace

import (
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc"
)

type TraceOption = otlptracegrpc.Option

type traceGrpcOption struct{}

func NewTraceGrpcOption() traceGrpcOption {
	return traceGrpcOption{}
}

func (tgo traceGrpcOption) WithGRPCConn(conn *grpc.ClientConn) TraceOption {
	return otlptracegrpc.WithGRPCConn(conn)
}

func (tgo traceGrpcOption) WithTimeout(duration time.Duration) TraceOption {
	return otlptracegrpc.WithTimeout(duration)
}
