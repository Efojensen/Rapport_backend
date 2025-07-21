package auth

import (
	"github.com/Efojensen/rapport.git/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")

	auth.Post("/student", handlers.StudentProfileSetup)

	auth.Post("/lecturer", handlers.LecturerProfileSetup)
}