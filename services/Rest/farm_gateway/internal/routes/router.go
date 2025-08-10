package routes

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	handlers []RouterHandlers
}

func NewRouter(handlers ...RouterHandlers) Router {
	return Router{
		handlers: handlers,
	}
}

func (r *Router) Builder(router fiber.Router) {
	for _, h := range r.handlers {
		router.Add(h.Method, h.Path, h.Handler...)
	}
}
