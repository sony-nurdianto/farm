package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/api"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/routes"
)

func main() {
	conn, err := api.NewGrpcClient("localhost:50051")
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	grpcSvc := api.NewGrpcService(conn.Connection())

	app := fiber.New()
	appRoutes := routes.NewRoutes(app, grpcSvc)
	appRoutes.Build()

	if err := app.Listen("0.0.0.0:3000"); err != nil {
		log.Fatalln(err)
	}
}
