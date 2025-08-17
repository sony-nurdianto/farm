package authh

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

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

	res, err := h.grpcAuthSvc.AuthUserSignIn(c.UserContext(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    res.Token,
		Expires:  res.ExpiresAt.AsTime(),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(
		fiber.Map{
			"data": fiber.Map{
				"status":    res.Status,
				"message":   res.Msg,
				"issued_at": res.IssuedAt.AsTime().Format(time.RFC3339),
			},
		},
	)
}
