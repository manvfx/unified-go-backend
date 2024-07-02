package database

import (
	"context"
	"time"
	"unified-go-backend/config"
	"unified-go-backend/utils"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectRedis(cfg *config.Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		utils.Logger.Fatalf("Failed to connect to Redis: %v", err)
	}

	utils.Logger.Info("Connected to Redis")
}

func DisconnectRedis() {
	if err := RedisClient.Close(); err != nil {
		utils.Logger.Fatalf("Failed to disconnect from Redis: %v", err)
	}

	utils.Logger.Info("Disconnected from Redis")
}
