package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/api"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/handlers"
)

type Routes struct {
	app     *fiber.App
	grpcSvc api.GrpcAuthService
}

func NewRoutes(app *fiber.App, grpcSvc api.GrpcAuthService) *Routes {
	return &Routes{
		app,
		grpcSvc,
	}
}

func (r *Routes) Build() {
	authHandler := handlers.NewAuthHandler(r.grpcSvc)

	signupHandler := NewRouterHandlers("/signup", http.MethodPost, authHandler.SignUp)
	signInHandler := NewRouterHandlers("/signin", http.MethodPost, authHandler.SignIn)
	authRouter := NewRouter(
		signupHandler,
		signInHandler,
	)
	r.app.Route("/auth", authRouter.Builder)
}
