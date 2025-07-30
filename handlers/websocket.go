package handlers

import (
	"log"

	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// HandleWebSocket manages WebSocket connections for chat
func HandleWebSocket(hub *models.Hub) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// Get user and room info from query parameters
		userId := c.Query("userId")
		roomId := c.Query("roomId")

		if userId == "" || roomId == "" {
			log.Printf("WebSocket connection rejected: missing userId or roomId")
			c.Close()
			return
		}

		// Create new client
		client := &models.Client{
			Conn:   c,
			UserId: userId,
			RoomId: roomId,
			Send:   make(chan models.WSMessage, 256),
		}

		// Register client with hub
		hub.Register <- client

		// Start goroutines for reading and writing
		go client.WritePump()
		go client.ReadPump(hub)
	})
}
