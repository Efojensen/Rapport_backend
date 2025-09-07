package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/Efojensen/rapport.git/constants"
	"github.com/Efojensen/rapport.git/handlers/secure"
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

	hash, err := secure.HashPassword(student.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	student.Password = hash

	result, err := collection.InsertOne(c.Context(), student)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": err.Error()})
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		student.ID = oid
	}

	res := fmt.Sprint(student.ID)

	communityUrl := fmt.Sprint(constants.PubUrl, res[10:34])

	proxy := fiber.AcquireAgent()
	proxy.Request().Header.SetMethod("POST")
	proxy.Request().SetRequestURI(communityUrl)

	err = proxy.Parse()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	var resBody struct {
		Message string `json:"message"`
	}

	_, body, errs := proxy.Bytes()
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(errs)
	}

	err = json.Unmarshal(body, &resBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg":          "student received successfully",
		"student":      student,
		"communityMsg": resBody.Message,
	})
}
