package handlers

import (
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
)

func SetupOther(c *fiber.Ctx) error {
	other := new(models.Other)

	if err := c.BodyParser(other); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": "Incorrect other body/fields",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg": "Other user created successfully",
		"other": other,
	})
}