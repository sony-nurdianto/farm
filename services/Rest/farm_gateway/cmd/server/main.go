package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/concurrent"
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

const (
	AuthServiceConnectionTag = iota
	FarmerServiceConnectionTag
	FarmServiceConnectionTag
)

type grpcConnClient struct {
	tag  int
	conn *grpc.ClientConn
}

func initGrpcClientConnection(
	ctx context.Context,
	tag int,
	target string,
	opts ...grpc.DialOption,
) <-chan any {
	out := make(chan any, 1)
	go func() {
		defer close(out)
		var res concurrent.Result[grpcConnClient]

		conn, err := grpc.NewClient(target, opts...)
		if err != nil {
			res.Error = err
			concurrent.SendResult(ctx, out, res)
			return
		}

		clientConn := grpcConnClient{
			tag:  tag,
			conn: conn,
		}

		res.Value = clientConn

		concurrent.SendResult(ctx, out, res)
	}()
	return out
}

func main() {
	godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serviceObsName := "farm-gateway"
	connColl, err := grpc.NewClient(
		os.Getenv("OTELCOLLECTORADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln(err)
	}

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

	statHandler := grpc.WithStatsHandler(
		otelgrpc.NewClientHandler(
			otelgrpc.WithTracerProvider(tp),
			otelgrpc.WithMeterProvider(mp),
			otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
		),
	)

	cred := grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	)

	clientIntercpt := api.NewUnaryClientInterceptor(tp)

	var authConnSvc *grpc.ClientConn
	var farmerConnSvc *grpc.ClientConn
	var farmConnSvc *grpc.ClientConn

	chs := []<-chan any{
		initGrpcClientConnection(
			ctx,
			AuthServiceConnectionTag,
			os.Getenv("GRPC_AUTH_SERVICE"),
			cred, statHandler,
			grpc.WithUnaryInterceptor(clientIntercpt.UnaryAuthClientIntercept),
		),
		initGrpcClientConnection(
			ctx,
			FarmerServiceConnectionTag,
			os.Getenv("GRPC_FARMER_SERVICE"),
			cred,
			statHandler,
			grpc.WithUnaryInterceptor(clientIntercpt.UnaryAuthClientIntercept),
		),
		initGrpcClientConnection(
			ctx,
			FarmServiceConnectionTag,
			os.Getenv("GRPC_FARM_SERVICE"),
			cred,
			statHandler,
			grpc.WithUnaryInterceptor(clientIntercpt.UnaryAuthClientIntercept),
		),
	}

	for v := range concurrent.FanIn(ctx, chs...) {
		res, ok := v.(concurrent.Result[grpcConnClient])
		if !ok {
			log.Fatalf("expected grpcClient conn but got %v\n", res)
		}

		if res.Error != nil {
			log.Fatalln(res.Error)
		}

		switch res.Value.tag {
		case AuthServiceConnectionTag:
			authConnSvc = res.Value.conn
		case FarmerServiceConnectionTag:
			farmerConnSvc = res.Value.conn
		case FarmServiceConnectionTag:
			farmConnSvc = res.Value.conn
		}

	}

	defer authConnSvc.Close()
	defer farmerConnSvc.Close()

	authSvc := api.NewGrpcService(pbgen.NewAuthServiceClient(authConnSvc))
	farmerSvc := api.NewGrpcFarmerService(pbgen.NewFarmerServiceClient(farmerConnSvc))
	farmSvc := api.NewGrpcFarmService(pbgen.NewFarmServiceClient(farmConnSvc))
	obsm := middleware.NewObservabilityMiddleware(tp, mp)

	app := fiber.New()
	app.Use(obsm.Trace)
	app.Use(obsm.Metric)
	appRoutes := routes.NewRoutes(app, authSvc, farmerSvc, farmSvc)
	appRoutes.Build()

	go func() {
		if err := app.Listen("0.0.0.0:3000"); err != nil {
			log.Fatalln(err)
		}
	}()

	var once sync.Once

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Server Stoping, Gracefully Stop ...")
			fmt.Println("Application Quit.")
			return
		default:
			once.Do(func() { fmt.Println("Server Running at 0.0.0.0:3000") })
		}
	}
}
