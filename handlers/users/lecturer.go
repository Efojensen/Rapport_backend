package handlers

import (
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func LecturerProfileSetup(c *fiber.Ctx, collection *mongo.Collection) error {
	lecturer := new(models.Lecturer)
	lecturer.Role = "lecturer"

	if err := c.BodyParser(lecturer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": "Invalid lecturer body",
		})
	}

	result, err := collection.InsertOne(c.Context(), lecturer)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		lecturer.ID = oid
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": "Failed to convert inserted ID to ObjectID",
		})
	}

	return c.JSON(fiber.Map{
		"msg":      "New lecturer created successfully",
		"lecturer": lecturer,
	})
}
