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
	// Note: You'll need to set proper SMTP credentials in your environment
	emailService := models.NewSMTPEmailService(
		"smtp.gmail.com", // host
		587,              // port
		sec_email,        // username (your email)
		"your-app-password", // password (use app-specific password for Gmail)
		"default@rapport.edu", // default recipient (not used for SOS emails)
	)

	// Setup routes
	email.SetupEmailRoutes(app, emailService, userCollection)
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
