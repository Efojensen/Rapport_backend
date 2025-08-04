package handlers

import (
	"github.com/Efojensen/rapport.git/handlers/secure"
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupOther(c *fiber.Ctx, collection *mongo.Collection) error {
	other := new(models.Other)
	other.Role = "other"

	if err := c.BodyParser(other); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": "Incorrect other body/fields",
		})
	}

	hash, err := secure.HashPassword(other.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	other.Password = hash

	_, err = collection.InsertOne(c.Context(), other)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg":   "Other user created successfully",
		"other": other,
	})
}
