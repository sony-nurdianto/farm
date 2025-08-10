package routes

import (
	"github.com/gofiber/fiber/v2"
)

type RouterHandlers struct {
	Path    string
	Method  string
	Handler []fiber.Handler
}

func NewRouterHandlers(path string, method string, handlers ...fiber.Handler) RouterHandlers {
	return RouterHandlers{
		Path:    path,
		Method:  method,
		Handler: handlers,
	}
}
