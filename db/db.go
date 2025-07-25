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

	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("User does not exist.")
			return nil, err
		}

		log.Println("Err: ", err)
		return nil, err
	}

	return user, nil
}
