package obs

import (
	"fmt"

	"go.opentelemetry.io/otel/metric"
)

type observeMeter struct {
	StartupDuration  metric.Float64Histogram
	DaemonUpTime     metric.Float64Counter
	RedisConnections metric.Int64UpDownCounter
	KafkaMessages    metric.Int64Counter
	ErrCounter       metric.Int64Counter
}

func InitObserveMeter(serviceName string, meter metric.Meter) observeMeter {
	startupDuration, _ := meter.Float64Histogram(
		fmt.Sprintf("%s_startup_duration_seconds", serviceName),
		metric.WithDescription("Time taken for service startup"),
		metric.WithUnit("s"),
	)

	daemonUpTime, _ := meter.Float64Counter(
		fmt.Sprintf("%s_daemon_uptime_seconds", serviceName),
		metric.WithDescription("Total uptime of the daemon in seconds"),
	)

	redisConnection, _ := meter.Int64UpDownCounter(
		fmt.Sprintf("%s_redis_connections", serviceName),
		metric.WithDescription("Number of active Redis connections"),
	)

	kafkaMessages, _ := meter.Int64Counter(
		fmt.Sprintf("%s_kafka_messages_total", serviceName),
		metric.WithDescription("Total number of Kafka messages processed"),
	)

	errorCounter, _ := meter.Int64Counter(
		fmt.Sprintf("%s_errors_total", serviceName),
		metric.WithDescription("Total number of errors encountered"),
	)

	return observeMeter{
		StartupDuration:  startupDuration,
		DaemonUpTime:     daemonUpTime,
		RedisConnections: redisConnection,
		KafkaMessages:    kafkaMessages,
		ErrCounter:       errorCounter,
	}
}
