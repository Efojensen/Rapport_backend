package main

import (
	"log"

	"github.com/Efojensen/rapport.git/db"
	"github.com/Efojensen/rapport.git/models"
	"github.com/Efojensen/rapport.git/routes/auth"
	"github.com/Efojensen/rapport.git/routes/email"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)


func main() {
	app := fiber.New()

	mongoClient, sec_email := db.ConnectToDb()
	userCollection := mongoClient.Database("Rapport").Collection("Users")

	emailService := models.NewSMTPEmailService(
		sec_email,
		587,
		"username",
		"password",
		"smtp.example.com",
	)

	// TODO Implement the logic to get the actual user id in the mongoDb database
	user, err := db.GetUserDetails("some-string", userCollection)
	if err != nil {
		log.Fatal("Something went wrong ")
	}

	email.SetupEmailRoutes(app, emailService, user)

	// TODO Implement websockets for chat feature
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

	auth.SetupAuthRoutes(app, userCollection)

	log.Fatal(app.Listen(":4000"))
}
