package farmh

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
)

func (fh farmHandler) CreateFarm(c *fiber.Ctx) error {
	farmerID := c.Locals("user_subject")
	id, ok := farmerID.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": errors.New("id is not string"),
			},
		)
	}

	var farm []models.CreateFarm

	if err := c.BodyParser(&farm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	for i := range len(farm) {
		farm[i].FarmerID = id
	}

	res, err := fh.grpcFarmSvc.CreateFarm(c.UserContext(), farm)
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
		"Msg":    "Create Farm Done",
	})
}
