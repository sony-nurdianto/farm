package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/repo"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/services"
	"github.com/sony-nurdianto/farm/services/Grpc/farm/internal/usescase"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/postgres/pkg"
	"github.com/sony-nurdianto/farm/shared_lib/Go/database/redis"
	"google.golang.org/grpc"
)

func main() {
	godotenv.Load()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	// serviceName := "farm-service"

	farmRepo, err := repo.NewFarmRepo(
		ctx, pkg.NewPostgresInstance(), redis.NewRedisInstance(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	farmUsecase := usescase.NewFarmUsecase(farmRepo)

	svc := services.NewFarmServiceServer(farmUsecase)
	svr := grpc.NewServer()

	pbgen.RegisterFarmServiceServer(svr, &svc)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	defer listener.Close()

	go func(lis net.Listener) {
		if err := svr.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}(listener)

	var once sync.Once

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Server Stoping, Gracefully Stop...")
			svr.GracefulStop()
			fmt.Println("Aplication Quit")
			return
		default:
			once.Do(func() { fmt.Println("Server Is Running At 0.0.0.0:50051") })

		}
	}
}
