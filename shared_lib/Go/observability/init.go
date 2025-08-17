package observability

import (
	"context"
	"time"

	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/metrics"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	logsdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

type observability struct {
	serviceName   string
	collectorConn *grpc.ClientConn
}

func NewObservability(
	serviceName string,
	collectorConn *grpc.ClientConn,
) observability {
	return observability{
		serviceName,
		collectorConn,
	}
}

func (o observability) Init(
	ctx context.Context,
) (*sdktrace.TracerProvider, *metric.MeterProvider, *logsdk.LoggerProvider, error) {
	var to time.Duration
	to = 15 * time.Second

	tcrEptr := trace.NewTraceGrpcExporter()
	tcrOpt := trace.NewTraceGrpcOption()
	tcr := trace.NewTracer(
		o.serviceName,
		tcrEptr,
	)

	tp, err := tcr.Provider(
		ctx,
		tcrOpt.WithGRPCConn(o.collectorConn),
		tcrOpt.WithTimeout(to),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	mtcEptr := metrics.NewMetricGrpcExporter()
	mtcOpt := metrics.NewMetricGrpcOption()
	mtc := metrics.NewMetric(
		o.serviceName,
		mtcEptr,
	)

	mp, err := mtc.Provider(
		ctx,
		mtcOpt.WithGRPCConn(o.collectorConn),
		mtcOpt.WithTimeout(to),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	lEptr := logs.NewLogGrpcExporter()
	lOpt := logs.NewLogGrpcOption()
	lo := logs.NewLogs(
		o.serviceName,
		lEptr,
	)

	lp, err := lo.Provider(
		ctx,
		lOpt.WithGRPCConn(o.collectorConn),
		lOpt.WithTimeout(to),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)
	global.SetLoggerProvider(lp)

	return tp, mp, lp, nil
}
