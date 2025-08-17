package authh

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

func (h authHandler) AuthTokenBaseValidate(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is missing",
		})
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Authorization header format",
		})
	}

	token := parts[1]

	req := &pbgen.TokenValidateRequest{
		Token: token,
	}

	res, err := h.grpcAuthSvc.AuthTokenValidate(c.UserContext(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	if !res.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": res.Msg,
			},
		)
	}

	c.Locals("user_subject", *res.Subject)
	c.Locals("user_isuer", *res.Isuer)
	c.Locals("user_experied", res.ExpiresAt.AsTime().Format(time.RFC3339))

	return c.Next()
}
