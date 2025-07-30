package email

import (
	"github.com/Efojensen/rapport.git/db"
	handlers "github.com/Efojensen/rapport.git/handlers/mail"
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupEmailRoutes(app *fiber.App, mailService *models.EmailService, userCollection *mongo.Collection) {
	mail := app.Group("/mail")

	// SOS email endpoint - expects userId in request body
	mail.Post("/sos", func(c *fiber.Ctx) error {
		// Parse the userId from request body
		var req struct {
			UserId string `json:"userId"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if req.UserId == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "userId is required",
			})
		}

		// Fetch user data from MongoDB
		user, err := db.GetUserDetails(req.UserId, userCollection)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Send the SOS email
		return handlers.SendGenericEmail(c, mailService, user)
	})
}
