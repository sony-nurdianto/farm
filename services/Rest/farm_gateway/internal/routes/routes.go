package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/api"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/handlers/authh"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/handlers/farmerh"
)

type Routes struct {
	app       *fiber.App
	authSvc   api.GrpcAuthService
	farmerSvc api.GrpcFarmerService
}

func NewRoutes(
	app *fiber.App,
	authSvc api.GrpcAuthService,
	farmerSvc api.GrpcFarmerService,
) *Routes {
	return &Routes{
		app,
		authSvc,
		farmerSvc,
	}
}

func (r *Routes) Build() {
	authHandler := authh.NewAuthHandler(r.authSvc)

	signupHandler := NewRouterHandlers("/signup", http.MethodPost, authHandler.SignUp)
	signInHandler := NewRouterHandlers("/signin", http.MethodPost, authHandler.SignIn)
	authRouter := NewRouter(
		signupHandler,
		signInHandler,
	)

	r.app.Route("/auth", authRouter.Builder)

	farmerHandler := farmerh.NewFarmerHandler(r.farmerSvc)
	farmerProfileHandler := NewRouterHandlers("/profile", http.MethodPost, authHandler.AuthTokenBaseValidate, farmerHandler.GetFarmerProfile)
	updateProfileHandler := NewRouterHandlers("/update_profile", http.MethodPatch, authHandler.AuthTokenBaseValidate, farmerHandler.UpdateUsers)
	farmerRouter := NewRouter(
		farmerProfileHandler,
		updateProfileHandler,
	)

	r.app.Route("/farmer", farmerRouter.Builder)

}
