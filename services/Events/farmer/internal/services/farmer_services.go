package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sony-nurdianto/farm/services/Events/farmer/internal/repo"
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
		"kafka_messages_processed_total",
		metric.WithDescription("Total number of Kafka messages processed"),
	)

	cacheOperations, _ := meter.Int64Counter(
		"cache_operations_total",
		metric.WithDescription("Total number of cache operations"),
	)

	msgCommitted, _ := meter.Int64Counter(
		"kafka_messages_committed_total",
		metric.WithDescription("Total number of kafka messages committed"),
	)

	prcsDuration, _ := meter.Float64Histogram(
		"message_processing_duration_seconds",
		metric.WithDescription("Time taken to process each message"),
		metric.WithUnit("s"),
	)

	deseDuration, _ := meter.Float64Histogram(
		"deserialization_duration_seconds",
		metric.WithDescription("Time taken to deserialize messages"),
		metric.WithUnit("s"),
	)

	cacheOperationDuration, _ := meter.Float64Histogram(
		"cache_operation_duration_seconds",
		metric.WithDescription("Time taken for cache operations (upsert/delete)"),
		metric.WithUnit("s"),
	)

	errorCounter, _ := meter.Int64Counter(
		"sync_errors_total",
		metric.WithDescription("Total number of errors during sync"),
	)

	activeConsumers, _ := meter.Int64UpDownCounter(
		"active_kafka_consumers",
		metric.WithDescription("Number of active Kafka consumers"),
	)

	fmCtx, span := tracer.Start(ctx, "sync_user_cache",
		trace.WithAttributes(
			attribute.String("kafka.topic", topic),
			attribute.String("operation", "sync_user_cache"),
		),
	)
	defer span.End()

	consumer := fs.repo.Consumer()
	consumer.SubscribeTopics([]string{topic}, kev.RebalanceCbCooperativeSticky)

	activeConsumers.Add(fmCtx, 1, metric.WithAttributes(
		attribute.String("topic", topic),
	))
	defer activeConsumers.Add(fmCtx, -1, metric.WithAttributes(
		attribute.String("topic", topic),
	))

	span.SetAttributes(attribute.String("consumer.status", "subscribed"))

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
				logger.Error(msgCtx, "Failed to read Kafka message", err)
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
			farmer, err := fs.repo.DeserializerFarmer(topic, msg.Value)
			deserDuration := time.Since(deserStart)

			deseDuration.Record(msgCtx, deserDuration.Seconds(),
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
				logger.Error(msgCtx, "Deserialization error", err)
				msgSpan.End()
				continue
			}

			msgSpan.SetAttributes(
				attribute.String("farmer.id", farmer.ID),
				attribute.Float64("deserialization.duration_seconds", deserDuration.Seconds()),
			)

			var op string
			for _, h := range msg.Headers {
				key := strings.TrimPrefix(h.Key, "__")
				if key == "op" {
					op = string(h.Value)
					break
				}
			}

			msgSpan.SetAttributes(attribute.String("cdc.operation", op))
			cacheKey := fmt.Sprintf("users:%s", farmer.ID)
			var cacheDuration time.Duration

			switch op {
			case "c", "u", "r":
				cacheStart := time.Now()
				cacheCtx, cacheSpan := tracer.Start(msgCtx, "cache_upsert",
					trace.WithAttributes(
						attribute.String("cache.key", cacheKey),
						attribute.String("farmer.id", farmer.ID),
						attribute.String("cache.operation", "upsert"),
					),
				)

				err := fs.repo.UpsertFarmerCache(cacheCtx, cacheKey, farmer)
				cacheDuration = time.Since(cacheStart)

				if err != nil {
					cacheSpan.SetStatus(codes.Error, "Cache upsert failed")
					cacheSpan.RecordError(err)
					errorCounter.Add(cacheCtx, 1, metric.WithAttributes(
						attribute.String("error.type", "cache_upsert_failed"),
						attribute.String("cache.key", cacheKey),
					))
					logger.Error(cacheCtx, "Cache upsert error", err)
					cacheSpan.End()
					msgSpan.End()
					continue
				}

				cacheSpan.SetStatus(codes.Ok, "Cache upsert successful")
				cacheOperations.Add(cacheCtx, 1, metric.WithAttributes(
					attribute.String("cache.operation", "upsert"),
					attribute.String("status", "success"),
				))

				cacheOperationDuration.Record(cacheCtx, cacheDuration.Seconds(),
					metric.WithAttributes(
						attribute.String("cache.operation", "upsert"),
					),
				)
				cacheSpan.End()

			case "d": // Delete
				cacheStart := time.Now()
				cacheCtx, cacheSpan := tracer.Start(msgCtx, "cache_delete",
					trace.WithAttributes(
						attribute.String("cache.key", cacheKey),
						attribute.String("farmer.id", farmer.ID),
						attribute.String("cache.operation", "delete"),
					),
				)

				err := fs.repo.DeleteFarmerCache(cacheCtx, cacheKey)
				cacheDuration = time.Since(cacheStart)

				if err != nil {
					cacheSpan.SetStatus(codes.Error, "Cache delete failed")
					cacheSpan.RecordError(err)
					errorCounter.Add(cacheCtx, 1, metric.WithAttributes(
						attribute.String("error.type", "cache_delete_failed"),
						attribute.String("cache.key", cacheKey),
					))
					logger.Error(cacheCtx, "Cache delete error", err)
					cacheSpan.End()
					msgSpan.End()
					continue
				}

				cacheSpan.SetStatus(codes.Ok, "Cache delete successful")
				cacheOperations.Add(cacheCtx, 1, metric.WithAttributes(
					attribute.String("cache.operation", "delete"),
					attribute.String("status", "success"),
				))

				// Fix: Record dengan operation yang benar
				cacheOperationDuration.Record(cacheCtx, cacheDuration.Seconds(),
					metric.WithAttributes(
						attribute.String("cache.operation", "delete"),
					),
				)
				cacheSpan.End()

			default:
				// Handle unknown operation
				msgSpan.SetAttributes(attribute.String("warning", "unknown_operation"))
				logger.Info(msgCtx, fmt.Sprintf("Unknown CDC operation: %s", op))
			}

			// Commit message
			commitStart := time.Now()
			if _, err := consumer.CommitMessage(msg); err != nil {
				msgSpan.SetStatus(codes.Error, "Message commit failed")
				msgSpan.RecordError(err)
				errorCounter.Add(msgCtx, 1, metric.WithAttributes(
					attribute.String("error.type", "kafka_commit_failed"),
					attribute.String("topic", topic),
				))
				logger.Error(msgCtx, "Kafka commit error", err)
				msgSpan.End()
				continue // Continue processing instead of return
			}

			commitDuration := time.Since(commitStart)
			msgCommitted.Add(msgCtx, 1, metric.WithAttributes(
				attribute.String("topic", topic),
				attribute.String("partition", fmt.Sprintf("%d", msg.TopicPartition.Partition)),
			))

			totalDuration := time.Since(startTime)
			prcsDuration.Record(msgCtx, totalDuration.Seconds(),
				metric.WithAttributes(
					attribute.String("topic", topic),
					attribute.String("operation", op),
					attribute.String("status", "success"),
				),
			)

			msgSpan.SetAttributes(
				attribute.Float64("processing.total_duration_seconds", totalDuration.Seconds()),
				attribute.Float64("processing.cache_duration_seconds", cacheDuration.Seconds()),
				attribute.Float64("processing.commit_duration_seconds", commitDuration.Seconds()),
				attribute.String("processing.status", "completed"),
			)
			msgSpan.SetStatus(codes.Ok, "Message processed successfully")
			msgSpan.End()
		}
	}
}
