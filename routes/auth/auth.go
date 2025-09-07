package auth

import (
	"github.com/Efojensen/rapport.git/handlers/users"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupAuthRoutes(app *fiber.App, authCollection *mongo.Collection) {
	auth := app.Group("/auth")

	auth.Get("/", func(c *fiber.Ctx) error {
		return handlers.GetAllUsers(c, authCollection)
	})

	auth.Get("/single", func (c *fiber.Ctx) error {
		return handlers.GetSingleUser(c, authCollection)
	})

	auth.Post("/register/student", func(c *fiber.Ctx) error {
		return handlers.StudentProfileSetup(c, authCollection)
	})

	auth.Post("/register/lecturer", func (c *fiber.Ctx) error  {
		return handlers.LecturerProfileSetup(c, authCollection)
	})

	auth.Post("/register/ta", func (c *fiber.Ctx) error {
		return handlers.SetupTeachAsst(c, authCollection)
	})

	auth.Post("/register/other", func(c *fiber.Ctx) error {
		return handlers.SetupOther(c, authCollection)
	})

	auth.Post("/login", func(c *fiber.Ctx) error {
		return handlers.UserLogin(c, authCollection)
	})
}