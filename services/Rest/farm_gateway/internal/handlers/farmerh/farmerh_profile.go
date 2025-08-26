package farmerh

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

func (h farmerHandler) GetFarmerProfile(c *fiber.Ctx) error {
	localID := c.Locals("user_subject")

	id, ok := localID.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": errors.New("id is not string"),
			},
		)
	}

	req := &pbgen.FarmerProfileRequest{
		Id: id,
	}

	res, err := h.grpcFarmerSvc.FarmerProfile(c.UserContext(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	user := models.Users{
		ID:           res.Farmer.Id,
		FullName:     res.Farmer.FullName,
		Email:        res.Farmer.Email,
		Phone:        res.Farmer.Phone,
		Verified:     res.Farmer.Verified,
		RegisteredAt: res.Farmer.RegisteredAt.AsTime().Format(time.RFC3339),
		UpdatedAt:    res.Farmer.UpdatedAt.AsTime().Format(time.RFC3339),
	}

	return c.JSON(user)
}
