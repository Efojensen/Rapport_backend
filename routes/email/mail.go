package email

import (
	"time"

	"github.com/Efojensen/rapport.git/db"
	handlers "github.com/Efojensen/rapport.git/handlers/mail"
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupEmailRoutes(app *fiber.App, mailService *models.EmailService, userCollection *mongo.Collection) {
	mail := app.Group("/mail")

	// SOS email endpoint - expects userId and optional location data in request body
	mail.Post("/sos", func(c *fiber.Ctx) error {
		// Parse the request body
		var req struct {
			UserId   string              `json:"userId"`
			Location *models.GeoLocation `json:"location,omitempty"`
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

		// Create SOS report with location if provided
		var sosReport *models.SOSReport
		if req.Location != nil {
			sosReport = &models.SOSReport{
				GeoLocation: *req.Location,
				SentTime:    time.Now(),
			}
		}

		// Send the SOS email with location data
		return handlers.SendSOSEmail(c, mailService, user, sosReport)
	})
}
