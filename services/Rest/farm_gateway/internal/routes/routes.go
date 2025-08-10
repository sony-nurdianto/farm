package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/handlers"
)

type Routes struct {
	app *fiber.App
}

func NewRoutes(app *fiber.App) *Routes {
	return &Routes{
		app: app,
	}
}

func (r *Routes) Build() {
	signupHandler := NewRouterHandlers("/signup", "GET", handlers.SignUp)
	signInHandler := NewRouterHandlers("/signin", "GET", handlers.SignIn)
	authRouter := NewRouter(
		signupHandler,
		signInHandler,
	)
	r.app.Route("/auth", authRouter.Builder)
}
