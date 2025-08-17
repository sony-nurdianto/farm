package logs

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
)

type LogsExporter = otlploggrpc.Exporter

type LogGrpcExporter interface {
	New(ctx context.Context, options ...LogsGrpcOptions) (*otlploggrpc.Exporter, error)
}

type logGrpcExporter struct{}

func NewLogGrpcExporter() logGrpcExporter {
	return logGrpcExporter{}
}

func (lge logGrpcExporter) New(ctx context.Context, options ...LogsGrpcOptions) (*LogsExporter, error) {
	return otlploggrpc.New(ctx, options...)
}
