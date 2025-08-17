package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/metrics"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

type observabilityMiddleware struct {
	traceProvider *trace.TraceProvider
	meterProvider *metrics.MeterProvider
	logger        *logs.Logger

	requestCount      metric.Int64Counter
	requestDuration   metric.Float64Histogram
	requestSize       metric.Int64Histogram
	responseSize      metric.Int64Histogram
	activeConnections metric.Int64UpDownCounter
	errorCount        metric.Int64Counter

	concurrentRequests metric.Int64UpDownCounter
	requestRate        metric.Float64Counter
}

func NewObservabilityMiddleware(
	tp *trace.TraceProvider,
	mp *metrics.MeterProvider,
) observabilityMiddleware {
	meter := mp.Meter("farm-gateway")

	requestCount, _ := meter.Int64Counter("http_server_request_count_total")
	requestDuration, _ := meter.Float64Histogram("http_server_request_duration_seconds", metric.WithUnit("s"))
	requestSize, _ := meter.Int64Histogram("http_server_request_size_bytes", metric.WithUnit("By"))
	responseSize, _ := meter.Int64Histogram("http_server_response_size_bytes", metric.WithUnit("By"))
	activeConnections, _ := meter.Int64UpDownCounter("http_server_active_connections")
	errorCount, _ := meter.Int64Counter("http_server_error_count_total")

	concurrentRequests, _ := meter.Int64UpDownCounter(
		"http_server_concurrent_requests",
		metric.WithDescription("Number of concurrent HTTP requests being processed"),
	)
	requestRate, _ := meter.Float64Counter(
		"http_server_request_rate",
		metric.WithDescription("HTTP request rate per second"),
	)

	return observabilityMiddleware{
		traceProvider: tp,
		meterProvider: mp,
		logger:        logs.NewLogger(),

		requestCount:      requestCount,
		requestDuration:   requestDuration,
		requestSize:       requestSize,
		responseSize:      responseSize,
		activeConnections: activeConnections,
		errorCount:        errorCount,

		concurrentRequests: concurrentRequests,
		requestRate:        requestRate,
	}
}

func (om observabilityMiddleware) Trace(c *fiber.Ctx) error {
	tp := om.traceProvider.Tracer("farm-gateway")

	ctx, span := tp.Start(
		c.UserContext(),
		c.Route().Path,
	)

	defer span.End()

	span.SetAttributes(
		attribute.String("http.method", c.Method()),
		attribute.String("http.route", c.Route().Path),
		attribute.String("http.url", c.OriginalURL()),
	)

	c.SetUserContext(ctx)

	if err := c.Next(); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("fiber.error", err.Error()))
		span.SetAttributes(attribute.Int("http.status_code", c.Response().StatusCode()))
		om.logger.Error(ctx, "fiber.error", err)
		return err
	}

	span.SetStatus(codes.Ok, "OK")
	span.SetAttributes(attribute.Int("http.status_code", c.Response().StatusCode()))

	return nil
}

func (om observabilityMiddleware) Metric(c *fiber.Ctx) error {
	start := time.Now()
	ctx := c.UserContext()

	// Standard attributes
	handler := c.Route().Name
	if handler == "" {
		handler = c.Route().Path
	}
	attrs := []attribute.KeyValue{
		attribute.String("http.method", c.Method()),
		attribute.String("http.route", c.Route().Path),
		attribute.String("http.handler", handler),
	}

	// Increment active connections / concurrent requests
	om.activeConnections.Add(ctx, 1, metric.WithAttributes(attrs...))
	defer om.activeConnections.Add(ctx, -1, metric.WithAttributes(attrs...))

	// Increment request rate
	om.requestRate.Add(ctx, 1.0, metric.WithAttributes(attrs...))

	// Calculate request size
	reqSize := len(c.Request().Header.Header()) + len(c.Body())

	// Process request
	err := c.Next()

	// Calculate duration and response size
	duration := time.Since(start).Seconds()
	respSize := len(c.Response().Body())

	statusCode := c.Response().StatusCode()
	statusClass := fmt.Sprintf("%dxx", statusCode/100)

	// Add dynamic attributes
	allAttrs := append(attrs,
		attribute.String("http.status_code", strconv.Itoa(statusCode)),
		attribute.String("http.status_class", statusClass),
	)
	if ua := c.Get("User-Agent"); ua != "" {
		allAttrs = append(allAttrs, attribute.String("http.user_agent", ua))
	}

	// Record core metrics
	om.requestCount.Add(ctx, 1, metric.WithAttributes(allAttrs...))
	om.requestDuration.Record(ctx, duration, metric.WithAttributes(allAttrs...))
	om.requestSize.Record(ctx, int64(reqSize), metric.WithAttributes(allAttrs...))
	om.responseSize.Record(ctx, int64(respSize), metric.WithAttributes(allAttrs...))

	// Record errors
	if statusCode >= 400 {
		errorAttrs := append(allAttrs, attribute.String("error_type", getErrorType(statusCode)))
		om.errorCount.Add(ctx, 1, metric.WithAttributes(errorAttrs...))
	}

	// Log error if any
	if err != nil {
		om.logger.Error(ctx, "fiber.error", err)
	}

	return err
}

// Helper function to categorize error types
func getErrorType(statusCode int) string {
	switch {
	case statusCode >= 400 && statusCode < 500:
		return "client_error"
	case statusCode >= 500:
		return "server_error"
	default:
		return "unknown"
	}
}
