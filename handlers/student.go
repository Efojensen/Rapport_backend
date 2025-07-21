package handlers

import (
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
)

func StudentProfileSetup(c *fiber.Ctx) error {
	student := new(models.Student)

	if err := c.BodyParser(student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid student body"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg": "student received successfully",
		"student": student,
	})
}