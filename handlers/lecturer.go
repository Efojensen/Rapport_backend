package handlers

import (
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
)

func LecturerProfileSetup(c *fiber.Ctx) error{
	lecturer := new(models.Lecturer)

	if err := c.BodyParser(lecturer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": "Invalid lecturer body",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "New lecturer created successfully",
		"lecturer": lecturer,
	})
}