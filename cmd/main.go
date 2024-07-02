package main

import (
	"context"
	"unified-go-backend/config"
	"unified-go-backend/database"
	"unified-go-backend/routes"
	"unified-go-backend/utils"
	"unified-go-backend/worker"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
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

	g, ctx := errgroup.WithContext(context.Background())

	// Start the email verification job worker
	g.Go(func() error {
		return worker.ProcessEmailVerificationJobs(ctx, cfg)
	})

	router := gin.Default()

	routes.AuthRoutes(router, cfg)
	routes.UserRoutes(router, cfg)

	utils.Logger.Info("Starting server...")

	g.Go(func() error {
		return router.Run()
	})

	if err := g.Wait(); err != nil {
		utils.Logger.Fatalf("Failed to run server: %v", err)
	}
}
