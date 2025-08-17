package logs

import (
	"context"

	logsdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type LoggerProvider = logsdk.LoggerProvider

type logs struct {
	svcName  string
	exporter LogGrpcExporter
}

func NewLogs(
	svcName string,
	exporter LogGrpcExporter,
) logs {
	return logs{
		svcName,
		exporter,
	}
}

func (l logs) Provider(
	ctx context.Context,
	opts ...LogsGrpcOptions,
) (*logsdk.LoggerProvider, error) {
	exporter, err := l.exporter.New(
		ctx,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	rsc := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(l.svcName),
	)

	lp := logsdk.NewLoggerProvider(
		logsdk.WithResource(rsc),
		logsdk.WithProcessor(
			logsdk.NewBatchProcessor(exporter),
		),
	)

	return lp, nil
}
