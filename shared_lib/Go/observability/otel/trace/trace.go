package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type TraceProvider = trace.TracerProvider

type Tracer interface {
	Provider(
		ctx context.Context,
		opts ...otlptracegrpc.Option,
	) (*trace.TracerProvider, error)
}

type tracer struct {
	svcName  string
	exporter TraceGrpcExporter
}

func NewTracer(
	svcName string,
	exporter TraceGrpcExporter,
) tracer {
	return tracer{
		svcName,
		exporter,
	}
}

func (t tracer) Provider(
	ctx context.Context,
	opts ...otlptracegrpc.Option,
) (*trace.TracerProvider, error) {
	exporter, err := t.exporter.New(
		ctx,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	rsc := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(t.svcName),
	)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(rsc),
	)

	return tp, nil
}
