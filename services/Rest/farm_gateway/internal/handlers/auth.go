package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/api"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

type authHandler struct {
	grpcAuthSvc api.GrpcAuthService
}

func NewAuthHandler(grpcSvc api.GrpcAuthService) authHandler {
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

	req := &pbgen.RegisterUserRequest{
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Password:    user.Password,
	}

	res, err := h.grpcAuthSvc.AuthUserRegister(req)
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
	var user models.UserSignIn

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	req := &pbgen.AuthenticateUserRequest{
		Email:    user.Email,
		Password: user.Password,
	}

	res, err := h.grpcAuthSvc.AuthUserSignIn(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(
		fiber.Map{
			"data": fiber.Map{
				"token":       res.Token,
				"status":      res.Status,
				"message":     res.Msg,
				"issued_at":   res.IssuedAt,
				"experied_at": res.ExpiresAt,
			},
		},
	)
}
