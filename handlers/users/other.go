package handlers

import (
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	result, err := collection.InsertOne(c.Context(), other)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		other.ID = oid
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": "Failed to convert inserted ID to ObjectID",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg": "Other user created successfully",
		"other": other,
	})
}