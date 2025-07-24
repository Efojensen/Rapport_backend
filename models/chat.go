package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RoomId    string             `json:"roomId" bson:"roomId"`
	SenderId  string             `json:"senderId" bson:"senderId"`
	Content   string             `json:"content" bson:"content"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	ReadBy    []string           `json:"readBy" bson:"readBy"`
}