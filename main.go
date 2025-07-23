package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Efojensen/rapport.git/routes/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	app := fiber.New()

	mongoClient := connectToDb()

	auth.SetupAuthRoutes(app, mongoClient)

	log.Fatal(app.Listen(":4000"))
}

func connectToDb() *mongo.Client {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoUrl := os.Getenv("MONGODB_URI")

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
	return client
}
