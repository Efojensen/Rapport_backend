package chats

import (
	"time"

	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateCommunityChat creates a college/department-wide community chat
func CreateCommunityChat(c *fiber.Ctx, roomCollection *mongo.Collection) error {
	var req struct {
		Name       string `json:"name"`
		CreatedBy  string `json:"createdBy"`
		College    string `json:"college,omitempty"`
		Department string `json:"department,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	room := models.Room{
		Id:        primitive.NewObjectID(),
		Name:      req.Name,
		Type:      models.CommunityChat,
		Members:   []string{}, // Community chats are open to all eligible users
		Admins:    []string{req.CreatedBy},
		CreatedBy: req.CreatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	result, err := roomCollection.InsertOne(c.Context(), room)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create community chat",
		})
	}

	room.Id = result.InsertedID.(primitive.ObjectID)
	return c.Status(fiber.StatusCreated).JSON(room)
}

// GetCommunityChats retrieves all community chats
func GetCommunityChats(c *fiber.Ctx, roomCollection *mongo.Collection) error {
	filter := bson.M{
		"type":     models.CommunityChat,
		"isActive": true,
	}

	cursor, err := roomCollection.Find(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve community chats",
		})
	}
	defer cursor.Close(c.Context())

	var rooms []models.Room
	if err = cursor.All(c.Context(), &rooms); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode community chats",
		})
	}

	return c.Status(fiber.StatusOK).JSON(rooms)
}

// JoinCommunityChat allows users to join a community chat
func JoinCommunityChat(c *fiber.Ctx, roomCollection *mongo.Collection) error {
	roomId := c.Params("roomId")
	userId := c.Query("userId")

	if roomId == "" || userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Room ID and User ID are required",
		})
	}

	objID, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	filter := bson.M{"_id": objID, "type": models.CommunityChat}
	update := bson.M{
		"$addToSet": bson.M{"members": userId},
		"$set":      bson.M{"updatedAt": time.Now()},
	}

	result, err := roomCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to join community chat",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Community chat not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully joined community chat",
	})
}
