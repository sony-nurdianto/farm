package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/api"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/routes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	godotenv.Load()

	insCred := insecure.NewCredentials()
	trsCred := grpc.WithTransportCredentials(insCred)

	conn, err := grpc.NewClient(os.Getenv("GRPC_AUTH_SERVICE"), trsCred)
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	authSvc := api.NewGrpcService(
		pbgen.NewAuthServiceClient(conn),
	)

	app := fiber.New()
	appRoutes := routes.NewRoutes(app, authSvc)
	appRoutes.Build()

	if err := app.Listen("0.0.0.0:3000"); err != nil {
		log.Fatalln(err)
	}
}
