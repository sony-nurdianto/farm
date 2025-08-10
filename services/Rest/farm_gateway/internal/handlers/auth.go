package handlers

import "github.com/gofiber/fiber/v2"

func SignUp(c *fiber.Ctx) error {
	return c.SendString("Hallo, World!")
}

func SignIn(c *fiber.Ctx) error {
	return c.SendString("SignIn")
}
