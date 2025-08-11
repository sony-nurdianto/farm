package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/api"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/handlers"
)

type Routes struct {
	app     *fiber.App
	grpcSvc api.GrpcService
}

func NewRoutes(app *fiber.App, grpcSvc api.GrpcService) *Routes {
	return &Routes{
		app,
		grpcSvc,
	}
}

func (r *Routes) Build() {
	authHandler := handlers.NewAuthHandler(r.grpcSvc.Services())

	signupHandler := NewRouterHandlers("/signup", "POST", authHandler.SignUp)
	signInHandler := NewRouterHandlers("/signin", "GET", authHandler.SignIn)
	authRouter := NewRouter(
		signupHandler,
		signInHandler,
	)
	r.app.Route("/auth", authRouter.Builder)
}
