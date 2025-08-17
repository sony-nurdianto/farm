package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type MeterProvider = sdkmetric.MeterProvider

type Metric interface {
	Provider(
		ctx context.Context,
		opts ...MetricOption,
	) (*sdkmetric.MeterProvider, error)
}

type metric struct {
	svcName  string
	exporter MetricGrpcExporter
}

func NewMetric(
	svcName string,
	exporter MetricGrpcExporter,
) metric {
	return metric{
		svcName,
		exporter,
	}
}

func (m metric) Provider(
	ctx context.Context,
	opts ...MetricOption,
) (*sdkmetric.MeterProvider, error) {
	exporter, err := m.exporter.New(
		ctx,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	rsc := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(m.svcName),
	)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(5*time.Second),
			),
		),
		sdkmetric.WithResource(rsc),
	)

	runtime.Start(
		runtime.WithMeterProvider(mp),
		runtime.WithMinimumReadMemStatsInterval(5*time.Second),
	)

	return mp, nil
}
