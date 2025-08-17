package logs

import (
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"google.golang.org/grpc"
)

type LogsGrpcOptions = otlploggrpc.Option

type logGrpcOption struct{}

func NewLogGrpcOption() logGrpcOption {
	return logGrpcOption{}
}

func (lgo logGrpcOption) WithGRPCConn(conn *grpc.ClientConn) LogsGrpcOptions {
	return otlploggrpc.WithGRPCConn(conn)
}

func (lgo logGrpcOption) WithTimeout(duration time.Duration) LogsGrpcOptions {
	return otlploggrpc.WithTimeout(duration)
}
