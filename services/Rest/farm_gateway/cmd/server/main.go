package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/api"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/middleware"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/routes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sony-nurdianto/farm/shared_lib/Go/observability"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
)

func main() {
	godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serviceObsName := "farm-gateway"
	connColl, err := grpc.NewClient(
		os.Getenv("OTELCOLLECTORADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	obs := observability.NewObservability(
		serviceObsName,
		connColl,
	)

	tp, mp, lp, err := obs.Init(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer tp.Shutdown(ctx)
	defer mp.Shutdown(ctx)
	defer lp.Shutdown(ctx)

	authCI := api.NewUnaryClientInterceptor(tp)

	conn, err := grpc.NewClient(
		os.Getenv("GRPC_AUTH_SERVICE"),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(
				otelgrpc.WithTracerProvider(tp),
				otelgrpc.WithMeterProvider(mp),
				otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
			),
		),

		grpc.WithUnaryInterceptor(authCI.UnaryAuthClientIntercept),
	)
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	authSvc := api.NewGrpcService(
		pbgen.NewAuthServiceClient(conn),
	)

	obsm := middleware.NewObservabilityMiddleware(tp, mp)

	app := fiber.New()
	app.Use(obsm.Trace)
	app.Use(obsm.Metric)
	appRoutes := routes.NewRoutes(app, authSvc)
	appRoutes.Build()

	go func() {
		if err := app.Listen("0.0.0.0:3000"); err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Server Stoping, Gracefully Stop ...")
			fmt.Println("Application Quit.")
			return
		default:
		}
	}
}
