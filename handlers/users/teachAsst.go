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

func SetupTeachAsst(c *fiber.Ctx, collection *mongo.Collection) error {
	ta := new(models.TeachAsst)
	ta.Role = "TA"

	if err := c.BodyParser(ta); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": "Invalid Teaching Assistant body",
		})
	}

	hash, err := secure.HashPassword(ta.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	ta.Password = hash

	result, err := collection.InsertOne(c.Context(), ta)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": err.Error()})
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		ta.ID = oid
	}

	res := fmt.Sprint(ta.ID)
	res = res[10:34]

	msg, err := utils.JoinCommunity(res)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg":          "student received successfully",
		"TA":           ta,
		"communityMsg": msg,
	})
}
