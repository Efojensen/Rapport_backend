package handlers

import (
	"strings"

	"github.com/Efojensen/rapport.git/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserLogin(c *fiber.Ctx, collection *mongo.Collection) error {
	var userCred struct {
		UsernameOrEmail string `json:"usernameOrEmail"`
		Password        string `json:"password"`
	}

	if err := c.BodyParser(&userCred); err != nil {
		return err
	}

	// Todo: Database call to find a user with matching email or password
	var err error;
	if strings.Contains(userCred.UsernameOrEmail, "@"){
		err = db.CheckUserCredByEmail(userCred.UsernameOrEmail, userCred.Password, collection)
	} else {
		err = db.CheckUserCredByUsername(userCred.UsernameOrEmail, userCred.Password, collection)
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": err,
		})
	}

	return nil
}
