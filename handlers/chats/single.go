package chats

import (
	"time"

	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateSingleChat creates or finds a single chat room between two users
func CreateSingleChat(c *fiber.Ctx, roomCollection *mongo.Collection) error {
	var req struct {
		User1 string `json:"user1"`
		User2 string `json:"user2"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if a single chat already exists between these users
	filter := bson.M{
		"type": models.SingleChat,
		"members": bson.M{
			"$all": []string{req.User1, req.User2},
			"$size": 2,
		},
		"isActive": true,
	}

	var existingRoom models.Room
	err := roomCollection.FindOne(c.Context(), filter).Decode(&existingRoom)
	if err == nil {
		// Room already exists
		return c.Status(fiber.StatusOK).JSON(existingRoom)
	}

	// Create new single chat room
	room := models.Room{
		Id:        primitive.NewObjectID(),
		Name:      "", // Single chats don't need names
		Type:      models.SingleChat,
		Members:   []string{req.User1, req.User2},
		Admins:    []string{}, // No admins in single chats
		CreatedBy: req.User1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	result, err := roomCollection.InsertOne(c.Context(), room)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create single chat",
		})
	}

	room.Id = result.InsertedID.(primitive.ObjectID)
	return c.Status(fiber.StatusCreated).JSON(room)
}

// SendChat sends a message to a chat room
func SendChat(c *fiber.Ctx, chatCollection *mongo.Collection, roomCollection *mongo.Collection, hub *models.Hub) error {
	chat := new(models.Chat)
	chat.Timestamp = time.Now()
	chat.MessageType = models.TextMessage // Default to text

	if err := c.BodyParser(chat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Problem parsing request body into chat",
		})
	}

	// Verify room exists and user is a member
	roomObjID, err := primitive.ObjectIDFromHex(chat.RoomId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	var room models.Room
	filter := bson.M{
		"_id": roomObjID,
		"members": bson.M{"$in": []string{chat.SenderId}},
		"isActive": true,
	}
	err = roomCollection.FindOne(c.Context(), filter).Decode(&room)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User not authorized to send messages to this room",
		})
	}

	// Insert the chat message
	result, err := chatCollection.InsertOne(c.Context(), chat)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send message",
		})
	}

	chat.Id = result.InsertedID.(primitive.ObjectID)

	// Update room's last message
	update := bson.M{
		"$set": bson.M{
			"lastMessage": chat,
			"updatedAt":   time.Now(),
		},
	}
	roomCollection.UpdateOne(c.Context(), bson.M{"_id": roomObjID}, update)

	// Broadcast message via WebSocket
	if hub != nil {
		wsMessage := models.WSMessage{
			Type:   "new_message",
			RoomId: chat.RoomId,
			UserId: chat.SenderId,
			Chat:   chat,
		}
		hub.Broadcast <- wsMessage
	}

	return c.Status(fiber.StatusCreated).JSON(chat)
}

// GetChatMessages retrieves messages for a specific room
func GetChatMessages(c *fiber.Ctx, chatCollection *mongo.Collection, roomCollection *mongo.Collection) error {
	roomId := c.Params("roomId")
	userId := c.Query("userId")
	limit := c.QueryInt("limit", 50) // Default 50 messages
	skip := c.QueryInt("skip", 0)   // For pagination

	if roomId == "" || userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Room ID and User ID are required",
		})
	}

	// Verify user is a member of the room
	roomObjID, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	var room models.Room
	filter := bson.M{
		"_id": roomObjID,
		"members": bson.M{"$in": []string{userId}},
		"isActive": true,
	}
	err = roomCollection.FindOne(c.Context(), filter).Decode(&room)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User not authorized to view messages in this room",
		})
	}

	// Get messages
	messageFilter := bson.M{"roomId": roomId}
	opts := options.Find().
		SetSort(bson.D{{"timestamp", -1}}). // Most recent first
		SetLimit(int64(limit)).
		SetSkip(int64(skip))

	cursor, err := chatCollection.Find(c.Context(), messageFilter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve messages",
		})
	}
	defer cursor.Close(c.Context())

	var messages []models.Chat
	if err = cursor.All(c.Context(), &messages); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode messages",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": messages,
		"room":     room,
	})
}

// GetSingleChats retrieves all single chats for a user
func GetSingleChats(c *fiber.Ctx, roomCollection *mongo.Collection) error {
	userId := c.Query("userId")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	filter := bson.M{
		"type":     models.SingleChat,
		"members":  bson.M{"$in": []string{userId}},
		"isActive": true,
	}

	cursor, err := roomCollection.Find(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve single chats",
		})
	}
	defer cursor.Close(c.Context())

	var rooms []models.Room
	if err = cursor.All(c.Context(), &rooms); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode single chats",
		})
	}

	return c.Status(fiber.StatusOK).JSON(rooms)
}