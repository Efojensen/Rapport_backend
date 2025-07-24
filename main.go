package main

import (
	"log"

	"github.com/Efojensen/rapport.git/db"
	"github.com/Efojensen/rapport.git/routes/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)


func main() {
	app := fiber.New()

	// Websocket
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
