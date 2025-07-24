package main

import (
	"log"

	"github.com/Efojensen/rapport.git/db"
	"github.com/Efojensen/rapport.git/models"
	"github.com/Efojensen/rapport.git/routes/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)


func main() {
	app := fiber.New()

	emailService := models.NewSMTPEmailService(
		"rapportsafety@gmail.com",
		587,
		"username",
		"password",
		"smtp.example.com",
	)

	app.Post("/send-email", func (c *fiber.Ctx) error {
		to := c.FormValue("to")
		subject := c.FormValue("subject")
		body := c.FormValue("body")

		if err := emailService.SendEmail(to, subject, body); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to send email",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Email sent successfully",
		})
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			if err := c.WriteMessage(mt, msg); err != nil {
				log.Println(msg)
			}
		}
	}))

	mongoClient := db.ConnectToDb()

	auth.SetupAuthRoutes(app, mongoClient)

	log.Fatal(app.Listen(":4000"))
}
