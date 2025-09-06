package handlers

import (
	"strings"

	"github.com/Efojensen/rapport.git/db"
	"github.com/Efojensen/rapport.git/handlers/jwt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserLogin(c *fiber.Ctx, collection *mongo.Collection) error {
	var userCred struct {
		UsernameOrEmail string `json:"usernameOrEmail"`
		Password        string `json:"password"`
	}

	if err := c.BodyParser(&userCred); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	// Todo: Database call to find a user with matching email or password
	var err error
	var uid string
	if strings.Contains(userCred.UsernameOrEmail, "@") {
		uid, err = db.CheckUserCredByEmail(userCred.UsernameOrEmail, userCred.Password, collection)
	} else {
		uid, err = db.CheckUserCredByUsername(userCred.UsernameOrEmail, userCred.Password, collection)
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	jwtToken, err := jwt.CreateJWT(userCred.UsernameOrEmail)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"_id": uid,
		"msg": "User login successful",
		"JWT": jwtToken,
	})
}
