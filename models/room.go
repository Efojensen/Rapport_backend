package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomType string

const (
	SingleChat    RoomType = "single"
	GroupChat     RoomType = "group"
	CommunityChat RoomType = "community"
)

type Room struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Type        RoomType           `json:"type" bson:"type"`
	Members     []string           `json:"members" bson:"members"`         // User IDs
	Admins      []string           `json:"admins" bson:"admins"`           // Admin User IDs (for groups)
	CreatedBy   string             `json:"createdBy" bson:"createdBy"`     // Creator User ID
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	LastMessage *Chat              `json:"lastMessage" bson:"lastMessage"` // For quick preview
	IsActive    bool               `json:"isActive" bson:"isActive"`
}

type RoomMember struct {
	UserId   string    `json:"userId" bson:"userId"`
	JoinedAt time.Time `json:"joinedAt" bson:"joinedAt"`
	Role     string    `json:"role" bson:"role"` // "member", "admin", "owner"
}
