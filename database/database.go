package database

import (
	"context"
	"time"
	"unified-go-backend/config"
	"unified-go-backend/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectDB(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		utils.Logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		utils.Logger.Fatalf("Failed to ping MongoDB: %v", err)
	}

	MongoClient = client
	utils.Logger.Info("Connected to MongoDB")
}

func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := MongoClient.Disconnect(ctx); err != nil {
		utils.Logger.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}

	utils.Logger.Info("Disconnected from MongoDB")
}
