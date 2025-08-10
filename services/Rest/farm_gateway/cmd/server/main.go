package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/routes"
)

func main() {
	app := fiber.New()

	appRoutes := routes.NewRoutes(app)
	appRoutes.Build()

	if err := app.Listen("0.0.0.0:3000"); err != nil {
		log.Fatalln(err)
	}
}
