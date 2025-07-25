package auth

import (
	"github.com/Efojensen/rapport.git/handlers/users"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupAuthRoutes(app *fiber.App, authCollection *mongo.Collection) {
	auth := app.Group("/auth")

	auth.Post("/student", func(c *fiber.Ctx) error {
		return handlers.StudentProfileSetup(c, authCollection)
	})

	auth.Post("/lecturer", func (c *fiber.Ctx) error  {
		return handlers.LecturerProfileSetup(c, authCollection)
	})

	auth.Post("/TA", func (c *fiber.Ctx) error {
		return handlers.SetupTeachAsst(c, authCollection)
	})

	auth.Post("/other", func(c *fiber.Ctx) error {
		return handlers.SetupOther(c, authCollection)
	})
}