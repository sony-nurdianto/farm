package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sony-nurdianto/farm/services/Events/farmer/internal/obs"
	"github.com/sony-nurdianto/farm/services/Events/farmer/internal/repo"
	"github.com/sony-nurdianto/farm/services/Events/farmer/internal/services"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sony-nurdianto/farm/shared_lib/Go/observability"
)

func redisConnection(
	ctx context.Context,
	tracer trace.Tracer,
	connMeter metric.Int64UpDownCounter,
	errMetter metric.Int64Counter,
) (c *redis.Client, err error) {
	rdCtx, span := tracer.Start(ctx, "redis_connection",
		trace.WithAttributes(
			attribute.String("component", "redis"),
			attribute.String("operation", "connect"),
			attribute.String("master_name", os.Getenv("FARMER_REDIS_MASTER_NAME")),
		),
	)

	defer span.End()

	count := 0
	for range 5 {
		attemptSpan := trace.SpanFromContext(rdCtx)
		attemptSpan.SetAttributes(attribute.Int("retry.attempt", count))

		rdb := redis.NewFailoverClient(
			&redis.FailoverOptions{
				MasterName: os.Getenv("FARMER_REDIS_MASTER_NAME"),
				SentinelAddrs: []string{
					os.Getenv("SENTINEL_FARMER_REDIS_ADDR"),
					os.Getenv("SENTINEL_FARMER_REDIS_ADDR_2"),
				},
				Username: os.Getenv("FARMER_REDIS_MASTER_USER_NAME"),
				Password: os.Getenv("FARMER_REDIS_MASTER_PASSWORD"),
				DB:       0,
			})

		_, err = rdb.Ping(context.Background()).Result()
		if err == nil {
			span.SetStatus(codes.Ok, "Redis connection successful")
			span.SetAttributes(
				attribute.Int("connection.attempts", count),
				attribute.String("connection.status", "success"),
			)

			connMeter.Add(rdCtx, 1)
			return rdb, nil
		}

		span.AddEvent("connection_attempt_failed",
			trace.WithAttributes(
				attribute.Int("attempt", count),
				attribute.String("error", err.Error()),
			),
		)

		time.Sleep(time.Second * 2)
		count++
	}

	span.SetStatus(codes.Error, "Redis connection failed after all retries")
	span.RecordError(err)
	errMetter.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("error.type", "redis_connection_failed"),
			attribute.String("component", "redis"),
		),
	)

	return nil, fmt.Errorf("connect failed after %d attempts: %w", count, err)
}

func main() {
	startTime := time.Now()

	godotenv.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serviceObsName := "auth-farmer-event"

	connColl, err := grpc.NewClient(
		os.Getenv("OTELCOLLECTORADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln(err)
	}

	obsProvider := observability.NewObservability(
		serviceObsName,
		connColl,
	)

	tp, mp, lp, err := obsProvider.Init(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer tp.Shutdown(ctx) // <- TraceProvider
	defer mp.Shutdown(ctx) // <- MeterProvider
	defer lp.Shutdown(ctx) // <-LogProvider

	tracer := tp.Tracer(serviceObsName)
	meter := mp.Meter(serviceObsName)
	logger := logs.NewLogger()

	observeMeter := obs.InitObserveMeter(
		serviceObsName,
		meter,
	)

	ctx, mainSpan := tracer.Start(ctx, "auth-service-daemon-lifecycle",
		trace.WithAttributes(
			attribute.String("service.name", serviceObsName),
		),
	)
	defer mainSpan.End()

	rdsCtx, rdSpan := tracer.Start(ctx, "initialize_redis_connection")
	rdb, err := redisConnection(rdsCtx, tracer, observeMeter.RedisConnections, observeMeter.ErrCounter)
	if err != nil {
		observeMeter.ErrCounter.Add(
			ctx, 1,
			metric.WithAttributes(
				attribute.String("service", serviceObsName),
				attribute.String("description", "error when init redis connection"),
				attribute.String("error_at", time.Now().Format(time.RFC3339)),
			),
		)

		rdSpan.SetStatus(codes.Error, "Redis initialization failed")
		rdSpan.RecordError(err)
		logger.Error(rdsCtx, err.Error(), err)
		return
	}
	rdSpan.SetStatus(codes.Ok, "Redis initialized successfully")
	rdSpan.End()

	cfgFarmer := map[kev.ConfigKeyKafka]string{
		kev.BOOTSTRAP_SERVERS:             os.Getenv("KAKFKABROKER"),
		kev.GROUP_ID:                      "farmer-event",
		kev.AUTO_OFFSET_RESET:             "earliest",
		kev.ENABLE_AUTO_COMMIT:            "false",
		kev.PARTITION_ASSIGNMENT_STRATEGY: "cooperative-sticky",
	}

	farmerSvcRepo, err := repo.NewFarmerRepo(
		schrgs.NewRegistery(),
		avr.NewAvrSerdeInstance(),
		kev.NewKafka(),
		cfgFarmer,
		pkg.NewPostgresInstance(),
		rdb,
	)
	if err != nil {
		observeMeter.ErrCounter.Add(
			ctx, 1,
			metric.WithAttributes(
				attribute.String("service", serviceObsName),
				attribute.String("description", "error when crate NewRepo"),
				attribute.String("error_at", time.Now().Format(time.RFC3339)),
			),
		)
		log.Fatalln(err)
	}

	defer farmerSvcRepo.CloseRepo()

	go func() {
		farmerSvc := services.NewFarmerService(farmerSvcRepo)
		if err := farmerSvc.SyncUserCache(
			ctx, "farmer-db.public.users_all_partitions",
			tracer,
			meter,
		); err != nil {
			observeMeter.ErrCounter.Add(
				ctx, 1,
				metric.WithAttributes(
					attribute.String("service", serviceObsName),
					attribute.String("description", "error when runninc SyncUserCache service"),
					attribute.String("error_at", time.Now().Format(time.RFC3339)),
				),
			)
			log.Fatalln(err)
		}
	}()

	observeMeter.StartupDuration.Record(
		ctx,
		time.Since(startTime).Seconds(),
		metric.WithAttributes(
			attribute.String("service", serviceObsName),
		),
	)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	uptimeStart := time.Now()

	var once sync.Once

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Event Server Stoping, Gracefully Stop ...")
			observeMeter.RedisConnections.Add(ctx, -1)
			fmt.Println("Application Quit.")
			return
		case <-ticker.C:
			observeMeter.DaemonUpTime.Add(
				ctx,
				time.Since(uptimeStart).Seconds(),
				metric.WithAttributes(
					attribute.String("service", serviceObsName),
				),
			)

		default:
			once.Do(func() { fmt.Println("event service daemon run") })
		}
	}
}
