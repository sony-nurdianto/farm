package metrics

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
)

type MetricExporter = otlpmetricgrpc.Exporter

type MetricGrpcExporter interface {
	New(ctx context.Context, options ...otlpmetricgrpc.Option) (*otlpmetricgrpc.Exporter, error)
}

type metricGrpcExporter struct{}

func NewMetricGrpcExporter() metricGrpcExporter {
	return metricGrpcExporter{}
}

func (mge metricGrpcExporter) New(ctx context.Context, options ...otlpmetricgrpc.Option) (*otlpmetricgrpc.Exporter, error) {
	return otlpmetricgrpc.New(ctx, options...)
}
