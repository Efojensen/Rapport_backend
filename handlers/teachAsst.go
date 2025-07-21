package handlers

import (
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
)

func SetupTeachAsst(c *fiber.Ctx) error {
	ta := new(models.TeachAsst)

	if err := c.BodyParser(ta); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": "Invalid Teaching Assistant body",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg": "New teaching assistant field created",
		"teachAsst": ta,
	})
}