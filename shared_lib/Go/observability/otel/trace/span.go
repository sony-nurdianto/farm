package trace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type traceSpan struct{}

func NewTraceSpan() traceSpan {
	return traceSpan{}
}

func (ts traceSpan) SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}
