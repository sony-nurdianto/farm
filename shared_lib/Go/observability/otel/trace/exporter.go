package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

type TraceExporter = otlptrace.Exporter

type TraceGrpcExporter interface {
	New(ctx context.Context, opts ...TraceOption) (*TraceExporter, error)
}

type traceGrpcExporter struct{}

func NewTraceGrpcExporter() traceGrpcExporter {
	return traceGrpcExporter{}
}

func (tge traceGrpcExporter) New(ctx context.Context, opts ...TraceOption) (*TraceExporter, error) {
	return otlptracegrpc.New(ctx, opts...)
}
