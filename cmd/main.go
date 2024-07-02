package main

import (
	"unified-go-backend/config"
	"unified-go-backend/database"
	"unified-go-backend/routes"
	"unified-go-backend/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize logger with Elasticsearch URL
	utils.InitLogger()

	// Connect to the database
	database.ConnectDB(cfg)
	defer database.DisconnectDB()

	// Connect to Redis
	database.ConnectRedis(cfg)
	defer database.DisconnectRedis()

	router := gin.Default()

	routes.AuthRoutes(router, cfg)
	routes.UserRoutes(router, cfg)

	utils.Logger.Info("Starting server...")

	if err := router.Run(); err != nil {
		utils.Logger.Fatalf("Failed to run server: %v", err)
	}
}
