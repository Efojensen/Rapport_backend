package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageType string

const (
	TextMessage   MessageType = "text"
	ImageMessage  MessageType = "image"
	FileMessage   MessageType = "file"
	SystemMessage MessageType = "system"
)

type Chat struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RoomId      string             `json:"roomId" bson:"roomId"`
	SenderId    string             `json:"senderId" bson:"senderId"`
	SenderName  string             `json:"senderName" bson:"senderName"`
	Content     string             `json:"content" bson:"content"`
	MessageType MessageType        `json:"messageType" bson:"messageType"`
	Timestamp   time.Time          `json:"timestamp" bson:"timestamp"`
	IsEdited    bool               `json:"isEdited" bson:"isEdited"`
	EditedAt    *time.Time         `json:"editedAt,omitempty" bson:"editedAt,omitempty"`
	ReplyTo     *string            `json:"replyTo,omitempty" bson:"replyTo,omitempty"` // Message ID being replied to
}