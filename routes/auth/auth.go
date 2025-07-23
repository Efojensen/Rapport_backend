package auth

import (
	"github.com/Efojensen/rapport.git/handlers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupAuthRoutes(app *fiber.App, client *mongo.Client) {
	authCollection := client.Database("Rapport").Collection("Users")
	auth := app.Group("/auth")

	auth.Post("/student", handlers.StudentProfileSetup)

	auth.Post("/lecturer", func (c *fiber.Ctx) error  {
		return handlers.LecturerProfileSetup(c, authCollection)
	})

	auth.Post("/TA", handlers.SetupTeachAsst)

	auth.Post("/other", handlers.SetupOther)
}