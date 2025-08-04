package handlers

import (
	"github.com/Efojensen/rapport.git/handlers/secure"
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func LecturerProfileSetup(c *fiber.Ctx, collection *mongo.Collection) error {
	lecturer := new(models.Lecturer)
	lecturer.Role = "lecturer"

	hash, err := secure.HashPassword(lecturer.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	lecturer.Password = hash

	if err := c.BodyParser(lecturer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": "Invalid lecturer body",
		})
	}

	_, err = collection.InsertOne(c.Context(), lecturer)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"msg":      "New lecturer created successfully",
		"lecturer": lecturer,
	})
}
