package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/interceptor"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/repo"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/service"
	"github.com/sony-nurdianto/farm/services/Grpc/farmer/internal/usecase"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability"
	"github.com/sony-nurdianto/farm/shared_lib/Go/observability/otel/logs"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	godotenv.Load()
	logger := logs.NewLogger()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serviceObsName := "farmer-service"

	connColl, err := grpc.NewClient(
		os.Getenv("OTELCOLLECTORADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal(ctx, err.Error(), err)
	}

	obs := observability.NewObservability(
		serviceObsName,
		connColl,
	)

	tp, mp, lp, err := obs.Init(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed Initiated Observability", err)
	}

	defer tp.Shutdown(ctx)
	defer mp.Shutdown(ctx)
	defer lp.Shutdown(ctx)

	svr := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(tp),
				otelgrpc.WithMeterProvider(mp),
				otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
			),
		),
		grpc.UnaryInterceptor(
			interceptor.AuthServiceUnaryInterceptor,
		),
	)

	farmerRepo, err := repo.NewFarmerRepo(
		ctx,
		pkg.NewPostgresInstance(),
		redis.NewRedisInstance(),
	)

	farmerUseCase := usecase.NewFarmerUseCase(farmerRepo)

	if err != nil {
		logger.Fatal(ctx, err.Error(), err)
	}

	svc := service.NewFarmerServiceServer(
		tp.Tracer(serviceObsName), mp.Meter(serviceObsName), farmerUseCase,
	)
	pbgen.RegisterFarmerServiceServer(svr, svc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal(ctx, "Failed To Create Listener", err)
	}

	go func(listener net.Listener) {
		if err := svr.Serve(listener); err != nil {
			logger.Fatal(context.Background(), err.Error(), err)
		}
	}(lis)

	var once sync.Once

	for {
		select {
		case <-ctx.Done():
			logger.Info(
				ctx,
				"Server Stoping, Gracefully Stop ...",
				slog.String("service_name", serviceObsName),
				slog.String("operations", "Gracefully Shutdown"),
			)
			svr.GracefulStop()
			logger.Info(
				ctx,
				"Aplication Quit ...",
				slog.String("service_name", serviceObsName),
				slog.String("operations", "Info"),
			)
			return
		default:
			once.Do(func() { fmt.Println("Server is Running ....") })
		}
	}
}
