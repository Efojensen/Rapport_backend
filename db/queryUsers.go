package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CheckUserCredByEmail(email, password string, collection *mongo.Collection) (string, error) {
	filter := bson.M{"email": email}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, get the role to determine which struct to use
	var roleDoc struct {
		Id       string `bson:"_id"`
		Role     string `bson:"role"`
		Password string `bson:"password"`
	}

	err := collection.FindOne(ctx, filter).Decode(&roleDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("user not found")
		}
		return "", err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(roleDoc.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return roleDoc.Id, nil
}

func CheckUserCredByUsername(username, password string, collection *mongo.Collection) (string, error) {
	filter := bson.M{"username": username}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, get the role to determine which struct to use
	var roleDoc struct {
		Id       string `bson:"_id"`
		Role     string `bson:"role"`
		Password string `bson:"password"`
	}

	err := collection.FindOne(ctx, filter).Decode(&roleDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("user not found")
		}
		return "", err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(roleDoc.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return roleDoc.Id, nil
}
