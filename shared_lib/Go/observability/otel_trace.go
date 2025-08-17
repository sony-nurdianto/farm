package observability

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type otelTrace struct{}

func NewOtelTrace() otelTrace {
	return otelTrace{}
}

func (ot otelTrace) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return otel.Tracer(name, opts...)
}

func (ot otelTrace) GetTracerProvider() trace.TracerProvider {
	return otel.GetTracerProvider()
}
