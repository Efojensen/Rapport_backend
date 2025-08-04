package handlers

import (
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupTeachAsst(c *fiber.Ctx, collection *mongo.Collection) error {
	ta := new(models.TeachAsst)
	ta.Role = "TA"

	if err := c.BodyParser(ta); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": "Invalid Teaching Assistant body",
		})
	}

	_, err := collection.InsertOne(c.Context(), ta)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg": "New teaching assistant field created",
		"teachAsst": ta,
	})
}