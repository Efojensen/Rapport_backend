package email

import (
	handlers "github.com/Efojensen/rapport.git/handlers/mail"
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
)

func SetupEmailRoutes(app *fiber.App, mailService *models.EmailService, user models.User) {
	mail := app.Group("/mail")

	mail.Post("", func (c *fiber.Ctx) error {
		return handlers.SendGenericEmail(c, mailService, user)
	})
}