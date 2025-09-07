package handlers

import (
	"github.com/Efojensen/rapport.git/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetSingleUser(c *fiber.Ctx, collection *mongo.Collection) error {
	userId := c.Params("userId")

	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID are required",
		})
	}

	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userId is required",
		})
	}

	// Fetch user data from MongoDB
	user, err := db.GetUserDetails(userId, collection)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
