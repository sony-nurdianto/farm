package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
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
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	repo, err := repository.NewAuthRepo(
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

	fmt.Println("Server Running At 0.0.0.0:50051")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Server Stoping, Gracefully Stop ...")
			gs.GracefulStop()
			fmt.Println("Application Quit.")
			return
		default:
		}
	}
}
