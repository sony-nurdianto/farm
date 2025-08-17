package recorderr

import (
	"context"
	"log/slog"
	"time"

	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	otelCodes "go.opentelemetry.io/otel/codes"
	otelTrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorRecorder interface {
	Record(
		ctx context.Context,
		code codes.Code,
		methdoName, msg string,
	) error
}

type errorRecorder struct {
	span  otelTrace.Span
	loggr *logs.Logger
}

func NewErrorRecorder(
	ts otelTrace.Span,
	loggr *logs.Logger,
) errorRecorder {
	return errorRecorder{ts, loggr}
}

func (ier errorRecorder) Record(
	ctx context.Context,
	code codes.Code,
	methdoName, msg string,
) error {
	slgAttr := []slog.Attr{
		slog.String("full_method", methdoName),
		slog.Int("codes", int(code)),
		slog.Time("error_at", time.Now()),
	}

	err := status.Error(code, msg)
	ier.span.RecordError(err)
	ier.span.SetStatus(otelCodes.Error, err.Error())
	ier.loggr.Error(ctx, msg, err, slgAttr...)

	return err
}
