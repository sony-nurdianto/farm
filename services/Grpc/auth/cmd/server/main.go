package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/codec"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/passencrypt"
	"github.com/sony-nurdianto/farm/auth/internal/encryption/token"
	"github.com/sony-nurdianto/farm/auth/internal/interceptor"
	"github.com/sony-nurdianto/farm/auth/internal/pbgen"
	"github.com/sony-nurdianto/farm/auth/internal/repository"
	"github.com/sony-nurdianto/farm/auth/internal/service"
	"github.com/sony-nurdianto/farm/auth/internal/usecase"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/avr"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/kev"
	"github.com/sony-nurdianto/farm/shared_lib/Go/kafkaev/schrgs"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sony-nurdianto/farm/shared_lib/Go/observability"
)

func main() {
	godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serviceObsName := "auth-service"

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

	repo, err := repository.NewAuthRepo(
		ctx,
		schrgs.NewRegistery(),
		pkg.NewPostgresInstance(),
		avr.NewAvrSerdeInstance(),
		kev.NewKafka(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	defer repo.CloseRepo()

	uc := usecase.NewServiceUsecase(
		&repo,
		passencrypt.NewPassEncrypt(
			rand.Reader,
			codec.NewBase64Encoder(),
		),
		token.NewTokhan(
			token.NewPassetoToken(),
		),
	)

	gs := grpc.NewServer(
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

	svc := service.NewAuthServiceServer(uc)
	pbgen.RegisterAuthServiceServer(gs, svc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	defer lis.Close()

	go func(listener net.Listener) {
		if err := gs.Serve(listener); err != nil {
			log.Fatalln(err)
		}
	}(lis)

	var once sync.Once

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Server Stoping, Gracefully Stop ...")
			gs.GracefulStop()
			fmt.Println("Application Quit.")
			return
		default:
			once.Do(func() {
				fmt.Println("Server Running At 0.0.0.0:50051")
			})
		}
	}
}
