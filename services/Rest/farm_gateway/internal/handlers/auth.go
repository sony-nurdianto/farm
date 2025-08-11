package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

type authHandler struct {
	grpcAuthSvc pbgen.AuthServiceClient
}

func NewAuthHandler(grpcSvc pbgen.AuthServiceClient) authHandler {
	return authHandler{
		grpcAuthSvc: grpcSvc,
	}
}

func (h authHandler) SignUp(c *fiber.Ctx) error {
	var user models.UserRegister

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	req := &pbgen.RegisterRequest{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Password:    user.Password,
	}

	res, err := h.grpcAuthSvc.Register(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"status": res.Status,
			"msg":    res.Msg,
		},
	)
}

func (h authHandler) SignIn(c *fiber.Ctx) error {
	return c.SendString("SignIn")
}
