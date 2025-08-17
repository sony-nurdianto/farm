package metrics

import (
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"google.golang.org/grpc"
)

type MetricOption = otlpmetricgrpc.Option

type metricGrpcOption struct{}

func NewMetricGrpcOption() metricGrpcOption {
	return metricGrpcOption{}
}

func (mgo metricGrpcOption) WithGRPCConn(conn *grpc.ClientConn) MetricOption {
	return otlpmetricgrpc.WithGRPCConn(conn)
}

func (mgo metricGrpcOption) WithTimeout(duration time.Duration) MetricOption {
	return otlpmetricgrpc.WithTimeout(duration)
}
