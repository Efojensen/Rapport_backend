package chat

import (
	"github.com/Efojensen/rapport.git/handlers/chats"
	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupChatRoutes configures all chat-related routes
func SetupChatRoutes(app *fiber.App, chatCollection *mongo.Collection, roomCollection *mongo.Collection, hub *models.Hub) {
	chatGroup := app.Group("/api/chat")

	// Single Chat Routes
	chatGroup.Post("/single/create", func(c *fiber.Ctx) error {
		return chats.CreateSingleChat(c, roomCollection)
	})
	chatGroup.Get("/single", func(c *fiber.Ctx) error {
		return chats.GetSingleChats(c, roomCollection)
	})

	// Group Chat Routes
	chatGroup.Post("/group/create", func(c *fiber.Ctx) error {
		return chats.CreateGroupChat(c, roomCollection)
	})
	chatGroup.Get("/group", func(c *fiber.Ctx) error {
		return chats.GetGroupChats(c, roomCollection)
	})
	chatGroup.Post("/group/:roomId/join", func(c *fiber.Ctx) error {
		return chats.JoinGroupChat(c, roomCollection)
	})
	chatGroup.Post("/group/:roomId/leave", func(c *fiber.Ctx) error {
		return chats.LeaveGroupChat(c, roomCollection)
	})

	// Community Chat Routes
	chatGroup.Post("/community/create", func(c *fiber.Ctx) error {
		return chats.CreateCommunityChat(c, roomCollection)
	})
	chatGroup.Get("/community", func(c *fiber.Ctx) error {
		return chats.GetCommunityChats(c, roomCollection)
	})
	chatGroup.Post("/community/:roomId/join", func(c *fiber.Ctx) error {
		return chats.JoinCommunityChat(c, roomCollection)
	})

	// Message Routes
	chatGroup.Post("/message/send", func(c *fiber.Ctx) error {
		return chats.SendChat(c, chatCollection, roomCollection, hub)
	})
	chatGroup.Get("/messages/:roomId", func(c *fiber.Ctx) error {
		return chats.GetChatMessages(c, chatCollection, roomCollection)
	})
}
