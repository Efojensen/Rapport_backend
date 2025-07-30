package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Efojensen/rapport.git/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDb() (*mongo.Client, string) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoUrl := os.Getenv("MONGODB_URI")
	sec_service := os.Getenv("SEC_SERVICE")

	clientOptions := options.Client().ApplyURI(mongoUrl)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatalf("Something went wrong: %s", err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	fmt.Printf("Connection to Database successful")
	return client, sec_service
}

func GetUserDetails(id string, collection *mongo.Collection) (models.User, error) {
	filter := bson.M{"_id": id}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, get the role to determine which struct to use
	var roleDoc struct {
		Role string `bson:"role"`
	}
	err := collection.FindOne(ctx, filter).Decode(&roleDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("User does not exist.")
			return nil, err
		}
		log.Println("Error getting user role: ", err)
		return nil, err
	}

	// Now decode into the appropriate struct based on role
	switch roleDoc.Role {
	case "student":
		var student models.Student
		err := collection.FindOne(ctx, filter).Decode(&student)
		if err != nil {
			log.Println("Error decoding student: ", err)
			return nil, err
		}
		return &student, nil
	case "TA":
		var ta models.TeachAsst
		err := collection.FindOne(ctx, filter).Decode(&ta)
		if err != nil {
			log.Println("Error decoding TA: ", err)
			return nil, err
		}
		return &ta, nil
	case "lecturer":
		var lecturer models.Lecturer
		err := collection.FindOne(ctx, filter).Decode(&lecturer)
		if err != nil {
			log.Println("Error decoding lecturer: ", err)
			return nil, err
		}
		return &lecturer, nil
	default:
		var other models.Other
		err := collection.FindOne(ctx, filter).Decode(&other)
		if err != nil {
			log.Println("Error decoding other user: ", err)
			return nil, err
		}
		return &other, nil
	}
}