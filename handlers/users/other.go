package handlers

import (
	"fmt"

	"github.com/Efojensen/rapport.git/handlers/secure"
	"github.com/Efojensen/rapport.git/models"
	"github.com/Efojensen/rapport.git/utils"
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

	hash, err := secure.HashPassword(other.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	other.Password = hash

	result, err := collection.InsertOne(c.Context(), other)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": err.Error()})
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		other.ID = oid
	}

	res := fmt.Sprint(other.ID)
	res = res[10:34]

	msg, err := utils.JoinCommunity(res)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg":          "student received successfully",
		"communityMsg": msg,
		"other":        other,
	})
}
