package handlers

import (
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func StudentProfileSetup(c *fiber.Ctx, collection *mongo.Collection) error {
	student := new(models.Student)
	student.Role = "student"

	if err := c.BodyParser(student); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid student body"})
	}

	result, err := collection.InsertOne(c.Context(), student)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": err})
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		student.ID = oid
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": "Failed to convert inserted ID to ObjectID",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg":     "student received successfully",
		"student": student,
	})
}
