package farmh

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
)

func (fh farmHandler) UpdateFarm(c *fiber.Ctx) error {
	var farmWitdhAddr []models.UpdateFarmWithAddr

	if err := c.BodyParser(&farmWitdhAddr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	res, err := fh.grpcFarmSvc.UpdateFarmOrAddress(c.UserContext(), farmWitdhAddr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(fiber.Map{
		"data":   res,
		"status": "Success",
		"msg":    "Update Farm Done",
	})
}
