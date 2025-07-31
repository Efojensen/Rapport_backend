package handlers

import (
	"context"
	"log"
	"time"

	"github.com/Efojensen/rapport.git/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllUsers(c *fiber.Ctx, collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch users",
			"data":    nil,
		})
	}
	defer cursor.Close(ctx)

	// Create a slice to hold all users
	var allUsers []interface{}

	// Iterate through cursor and decode each user based on their role
	for cursor.Next(ctx) {
		// Peek at the role and decode into the correct struct
		var rawDoc bson.M
		if err := cursor.Decode(&rawDoc); err != nil {
			log.Printf("Error decoding document: %v", err)
			continue
		}

		role, _ := rawDoc["role"].(string)
		switch role {
		case "student":
			var student models.Student
			bsonBytes, _ := bson.Marshal(rawDoc)
			if err := bson.Unmarshal(bsonBytes, &student); err != nil {
				log.Printf("Error decoding student: %v", err)
				continue
			}
			allUsers = append(allUsers, student)
		case "TA":
			var ta models.TeachAsst
			bsonBytes, _ := bson.Marshal(rawDoc)
			if err := bson.Unmarshal(bsonBytes, &ta); err != nil {
				log.Printf("Error decoding TA: %v", err)
				continue
			}
			allUsers = append(allUsers, ta)
		case "lecturer":
			var lecturer models.Lecturer
			bsonBytes, _ := bson.Marshal(rawDoc)
			if err := bson.Unmarshal(bsonBytes, &lecturer); err != nil {
				log.Printf("Error decoding lecturer: %v", err)
				continue
			}
			allUsers = append(allUsers, lecturer)
		default:
			var other models.Other
			bsonBytes, _ := bson.Marshal(rawDoc)
			if err := bson.Unmarshal(bsonBytes, &other); err != nil {
				log.Printf("Error decoding other user: %v", err)
				continue
			}
			allUsers = append(allUsers, other)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error iterating through users",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Users retrieved successfully",
		"count":   len(allUsers),
		"data":    allUsers,
	})
}
