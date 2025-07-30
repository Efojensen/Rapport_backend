package chats

import (
	"time"

	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateGroupChat creates a new group chat room
func CreateGroupChat(c *fiber.Ctx, roomCollection *mongo.Collection) error {
	var req struct {
		Name      string   `json:"name"`
		Members   []string `json:"members"`
		CreatedBy string   `json:"createdBy"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Add creator to members if not already included
	found := false
	for _, member := range req.Members {
		if member == req.CreatedBy {
			found = true
			break
		}
	}
	if !found {
		req.Members = append(req.Members, req.CreatedBy)
	}

	room := models.Room{
		Id:        primitive.NewObjectID(),
		Name:      req.Name,
		Type:      models.GroupChat,
		Members:   req.Members,
		Admins:    []string{req.CreatedBy}, // Creator is admin
		CreatedBy: req.CreatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	result, err := roomCollection.InsertOne(c.Context(), room)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create group chat",
		})
	}

	room.Id = result.InsertedID.(primitive.ObjectID)
	return c.Status(fiber.StatusCreated).JSON(room)
}

// JoinGroupChat adds a user to an existing group chat
func JoinGroupChat(c *fiber.Ctx, roomCollection *mongo.Collection) error {
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

	filter := bson.M{"_id": objID, "type": models.GroupChat}
	update := bson.M{
		"$addToSet": bson.M{"members": userId},
		"$set":      bson.M{"updatedAt": time.Now()},
	}

	result, err := roomCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to join group chat",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Group chat not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully joined group chat",
	})
}

// LeaveGroupChat removes a user from a group chat
func LeaveGroupChat(c *fiber.Ctx, roomCollection *mongo.Collection) error {
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

	filter := bson.M{"_id": objID, "type": models.GroupChat}
	update := bson.M{
		"$pull": bson.M{"members": userId, "admins": userId},
		"$set":  bson.M{"updatedAt": time.Now()},
	}

	result, err := roomCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to leave group chat",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Group chat not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully left group chat",
	})
}

// GetGroupChats retrieves all group chats for a user
func GetGroupChats(c *fiber.Ctx, roomCollection *mongo.Collection) error {
	userId := c.Query("userId")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	filter := bson.M{
		"type":    models.GroupChat,
		"members": bson.M{"$in": []string{userId}},
		"isActive": true,
	}

	cursor, err := roomCollection.Find(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve group chats",
		})
	}
	defer cursor.Close(c.Context())

	var rooms []models.Room
	if err = cursor.All(c.Context(), &rooms); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode group chats",
		})
	}

	return c.Status(fiber.StatusOK).JSON(rooms)
}