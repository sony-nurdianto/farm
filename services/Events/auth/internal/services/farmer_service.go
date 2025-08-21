package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sony-nurdianto/farm/services/Events/auth/internal/repo"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type farmerService struct {
	repo repo.Repo
}

func NewFarmerService(rp repo.Repo) farmerService {
	return farmerService{
		repo: rp,
	}
}

func (fs farmerService) SyncUserCache(
	ctx context.Context,
	topic string,
	tracer trace.Tracer,
	meter metric.Meter,
) error {
	logger := logs.NewLogger()

	msgProcessed, _ := meter.Int64Counter(
		"kafka_messages_process_total",
		metric.WithDescription("Total number of Kafka messages processed"),
	)

	cacheOperations, _ := meter.Int64Counter(
		"cache_operations_total",
		metric.WithDescription("Total number of cache operations"),
	)

	msgCommited, _ := meter.Int64Counter(
		"kafka_messages_commited_total",
		metric.WithDescription("Total number of kafka messages commited"),
	)

	prcsDuration, _ := meter.Float64Histogram(
		"messages_processing_duration_seconds",
		metric.WithDescription("Time taken to process each message"),
		metric.WithUnit("s"),
	)

	deseDuratoin, _ := meter.Float64Histogram(
		"deserialization_duration_seconds",
		metric.WithDescription("Time taken to deserialize messages"),
		metric.WithUnit("s"),
	)

	insrtCacheDuration, _ := meter.Float64Histogram(
		"cache_insert_duration_seconds",
		metric.WithDescription("Time taken to insert into cache"),
		metric.WithUnit("s"),
	)

	errorCounter, _ := meter.Int64Counter(
		"sync_errors_total",
		metric.WithDescription("Total number of errors during sync"),
	)

	fmCtx, span := tracer.Start(ctx, "sync_user_cache",
		trace.WithAttributes(
			attribute.String("kafka.topic", topic),
			attribute.String("operation", "sync_user_cache"),
		),
	)

	defer span.End()

	span.SetAttributes(attribute.String("consumer.status", "subscribed"))
	consumer := fs.repo.Consumer()
	consumer.SubscribeTopics([]string{topic}, kev.RebalanceCbCooperativeSticky)

	for {
		select {
		case <-ctx.Done():
			span.SetStatus(codes.Ok, "Sync stopped gracefully")
			return ctx.Err()
		default:
			msgCtx, msgSpan := tracer.Start(fmCtx, "process_kafka_message")
			startTime := time.Now()
			msg, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if _, ok := err.(kev.KevError); ok {
					msgSpan.SetStatus(codes.Ok, "No message available")
					msgSpan.End()
					continue
				}

				msgSpan.SetStatus(codes.Error, "Failed to read message")
				msgSpan.RecordError(err)
				errorCounter.Add(msgCtx, 1, metric.WithAttributes(
					attribute.String("error.type", "kafka_read_failed"),
					attribute.String("topic", topic),
				))
				msgSpan.End()
				return err

			}

			msgSpan.SetAttributes(
				attribute.String("kafka.partition", fmt.Sprintf("%d", msg.TopicPartition.Partition)),
				attribute.String("kafka.offset", fmt.Sprintf("%d", msg.TopicPartition.Offset)),
				attribute.Int("message.size", len(msg.Value)),
			)

			msgProcessed.Add(msgCtx, 1, metric.WithAttributes(
				attribute.String("topic", topic),
				attribute.String("partition", fmt.Sprintf("%d", msg.TopicPartition.Partition)),
			))

			deserStart := time.Now()
			// call MethodDeserializer
			farmer, err := fs.repo.DeserializerFarmer(topic, msg.Value)
			deserDuration := time.Since(deserStart)

			deseDuratoin.Record(msgCtx, deserDuration.Seconds(),
				metric.WithAttributes(
					attribute.String("topic", topic),
				),
			)

			if err != nil {
				msgSpan.SetStatus(codes.Error, "Deserialization failed")
				msgSpan.RecordError(err)
				errorCounter.Add(msgCtx, 1, metric.WithAttributes(
					attribute.String("error.type", "deserialization_failed"),
					attribute.String("topic", topic),
				))

				logger.Error(msgCtx, fmt.Sprintf("Deserialization error: %s", err.Error()), err)
				msgSpan.End()
				continue
			}

			msgSpan.SetAttributes(
				attribute.String("farmer.id", farmer.ID),
				attribute.Float64("deserialization.duration_seconds", deserDuration.Seconds()),
			)

			cacheKey := fmt.Sprintf("users:%s", farmer.ID)
			cacheStart := time.Now()

			cacheCtx, cacheSpan := tracer.Start(msgCtx, "cache_insert",
				trace.WithAttributes(
					attribute.String("cache.key", cacheKey),
					attribute.String("farmer.id", farmer.ID),
				),
			)

			// Call Method InsertFarmerCache
			if err := fs.repo.InsertFarmerCache(cacheCtx, cacheKey, farmer); err != nil {
				cacheSpan.SetStatus(codes.Error, "Cache insertion failed")
				cacheSpan.RecordError(err)
				errorCounter.Add(cacheCtx, 1, metric.WithAttributes(
					attribute.String("error.type", "cache_insert_failed"),
					attribute.String("cache.key", cacheKey),
				))

				logger.Error(cacheCtx, fmt.Sprintf("Cache insertion error:%s", err.Error()), err)
				continue
			}

			cacheSpan.SetStatus(codes.Ok, "Cache inserted successfully")
			cacheOperations.Add(cacheCtx, 1, metric.WithAttributes(
				attribute.String("cache.operation", "insert"),
				attribute.String("status", "success"),
			))

			cacheDuration := time.Since(cacheStart)
			insrtCacheDuration.Record(cacheCtx, cacheDuration.Seconds(),
				metric.WithAttributes(
					attribute.String("cache.operation", "insert"),
				),
			)

			cacheSpan.End()

			commitStart := time.Now()
			// Call Method CommitMessage
			if _, err := consumer.CommitMessage(msg); err != nil {
				msgSpan.SetStatus(codes.Error, "Message commit failed")
				msgSpan.RecordError(err)
				errorCounter.Add(msgCtx, 1, metric.WithAttributes(
					attribute.String("error.type", "kafka_commit_failed"),
					attribute.String("topic", topic),
				))
				msgSpan.End()
				logger.Error(msgCtx, err.Error(), err)
				continue
			}

			msgCommited.Add(msgCtx, 1, metric.WithAttributes(
				attribute.String("topic", topic),
				attribute.String("partition", fmt.Sprintf("%d", msg.TopicPartition.Partition)),
			))

			totalDuration := time.Since(startTime)
			prcsDuration.Record(msgCtx, totalDuration.Seconds(),
				metric.WithAttributes(
					attribute.String("topic", topic),
					attribute.String("status", "success"),
				),
			)

			msgSpan.SetAttributes(
				attribute.Float64("processing.total_duration_seconds", totalDuration.Seconds()),
				attribute.Float64("processing.cache_duration_seconds", cacheDuration.Seconds()),
				attribute.Float64("processing.commit_duration_seconds", time.Since(commitStart).Seconds()),
				attribute.String("processing.status", "completed"),
			)
			msgSpan.SetStatus(codes.Ok, "Message processed successfully")
			msgSpan.End()
		}
	}
}
