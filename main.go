package main

import (
	"log"

	"github.com/Efojensen/rapport.git/db"
	"github.com/Efojensen/rapport.git/handlers"
	"github.com/Efojensen/rapport.git/models"
	"github.com/Efojensen/rapport.git/routes/auth"
	"github.com/Efojensen/rapport.git/routes/chat"
	"github.com/Efojensen/rapport.git/routes/email"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Enable CORS for frontend integration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Database connections
	mongoClient, sec_email := db.ConnectToDb()
	userCollection := mongoClient.Database("Rapport").Collection("Users")
	chatCollection := mongoClient.Database("Rapport").Collection("Chats")
	roomCollection := mongoClient.Database("Rapport").Collection("Rooms")

	// Initialize WebSocket hub for real-time chat
	hub := models.NewHub()
	go hub.Run()

	// Email service setup
	emailService := models.NewSMTPEmailService(
		sec_email,
		587,
		"username",
		"password",
		"smtp.gmail.com",
	)

	// TODO Implement the logic to get the actual user id in the mongoDb database
	user, err := db.GetUserDetails("some-string", userCollection)
	if err != nil {
		log.Printf("Warning: Could not get user details: %v", err)
		// Don't fatal here, continue with other services
	}

	// Setup routes
	email.SetupEmailRoutes(app, emailService, user)
	auth.SetupAuthRoutes(app, userCollection)
	chat.SetupChatRoutes(app, chatCollection, roomCollection, hub)

	// WebSocket endpoint for real-time chat
	app.Get("/ws", handlers.HandleWebSocket(hub))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Rapport API is running",
		})
	})

	log.Println("🚀 Rapport server starting on port 4000")
	log.Println("📱 WebSocket endpoint: ws://localhost:4000/ws")
	log.Println("🔗 API endpoints: http://localhost:4000")
	log.Fatal(app.Listen(":4000"))
}
