package farmerh

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

func (h farmerHandler) GetFarmerProfile(c *fiber.Ctx) error {

	dr := struct {
		ID string `json:"id"`
	}{}

	err := c.BodyParser(&dr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	req := &pbgen.FarmerProfileRequest{
		Id: dr.ID,
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
