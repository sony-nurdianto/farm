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

func (l *Logger) logToOtel(ctx context.Context, msg string, severity otellog.Severity, attrs []slog.Attr) {
	record := otellog.Record{}
	record.SetBody(otellog.StringValue(msg))
	record.SetSeverity(severity)
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

func (l *Logger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.slogLogger.InfoContext(ctx, msg, attrsToAny(attrs)...)
	l.logToOtel(ctx, msg, otellog.SeverityInfo, attrs)
}

func (l *Logger) Error(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	var allAttrs []slog.Attr
	if err != nil {
		allAttrs = append(attrs, slog.String("error", err.Error()))
	} else {
		allAttrs = attrs
	}

	l.slogLogger.ErrorContext(ctx, msg, attrsToAny(allAttrs)...)
	l.logToOtel(ctx, msg, otellog.SeverityError, allAttrs)
}

func (l *Logger) Fatal(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	var allAttrs []slog.Attr
	if err != nil {
		allAttrs = append(attrs, slog.String("error", err.Error()))
	} else {
		allAttrs = attrs
	}

	l.slogLogger.ErrorContext(ctx, msg, attrsToAny(allAttrs)...)
	l.logToOtel(ctx, msg, otellog.SeverityFatal, allAttrs)

	os.Exit(1)
}

func attrsToAny(attrs []slog.Attr) []any {
	result := make([]any, 0, len(attrs)*2)
	for _, attr := range attrs {
		result = append(result, attr.Key, attr.Value)
	}
	return result
}
