package farmerh

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/pbgen"
)

func (h farmerHandler) UpdateUsers(c *fiber.Ctx) error {
	var user models.UpdateUsers

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	req := &pbgen.UpdateFarmerProfileRequest{
		Id:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Phone:    user.Phone,
	}

	res, err := h.grpcFarmerSvc.ProfileFarmerUpdate(c.UserContext(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(fiber.Map{
		"status": res.Status,
		"msg":    res.Msg,
	})
}
