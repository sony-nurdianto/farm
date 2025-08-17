package logs

import (
	"context"
	"log/slog"
	"os"
	"time"

	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	otelLogger otellog.Logger
	slogLogger *slog.Logger
}

func NewLogger() *Logger {
	return &Logger{
		otelLogger: global.GetLoggerProvider().Logger("my-service"),
		slogLogger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	}
}

func (l *Logger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	// Log to slog (stdout)
	l.slogLogger.InfoContext(ctx, msg, attrsToAny(attrs)...)

	// Log to OpenTelemetry
	record := otellog.Record{}
	record.SetBody(otellog.StringValue(msg))
	record.SetSeverity(otellog.SeverityInfo)
	record.SetTimestamp(time.Now())

	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		sc := span.SpanContext()
		record.AddAttributes(
			otellog.String("trace_id", sc.TraceID().String()),
			otellog.String("span_id", sc.SpanID().String()),
		)
	}

	for _, attr := range attrs {
		record.AddAttributes(otellog.String(attr.Key, attr.Value.String()))
	}

	l.otelLogger.Emit(ctx, record)
}

func (l *Logger) Error(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	allAttrs := append(attrs, slog.String("error", err.Error()))

	l.slogLogger.ErrorContext(ctx, msg, attrsToAny(allAttrs)...)

	// Log to OpenTelemetry
	record := otellog.Record{}
	record.SetBody(otellog.StringValue(msg))
	record.SetSeverity(otellog.SeverityError)
	record.SetTimestamp(time.Now())

	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		sc := span.SpanContext()
		record.AddAttributes(
			otellog.String("trace_id", sc.TraceID().String()),
			otellog.String("span_id", sc.SpanID().String()),
		)
	}

	for _, attr := range allAttrs {
		record.AddAttributes(otellog.String(attr.Key, attr.Value.String()))
	}

	l.otelLogger.Emit(ctx, record)
}

func attrsToAny(attrs []slog.Attr) []any {
	result := make([]any, 0, len(attrs)*2)
	for _, attr := range attrs {
		result = append(result, attr.Key, attr.Value)
	}
	return result
}
