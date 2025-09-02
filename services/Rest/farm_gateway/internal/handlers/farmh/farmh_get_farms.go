package farmh

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sony-nurdianto/farm/services/Rest/farm_gateway/farm_gateway/internal/models"
)

func (fh farmHandler) GetFarmByID(c *fiber.Ctx) error {
	req := struct {
		ID string `json:"id"`
	}{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	res, err := fh.grpcFarmSvc.GetFarmByID(c.UserContext(), req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": errors.New("id is not string"),
			},
		)
	}

	return c.JSON(res)
}

func (fh farmHandler) GetFarms(c *fiber.Ctx) error {
	farmerID := c.Locals("user_subject")
	id, ok := farmerID.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": errors.New("id is not string"),
			},
		)
	}

	var req models.GetFarmsRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	res, err := fh.grpcFarmSvc.GetFarms(c.UserContext(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(res)
}
