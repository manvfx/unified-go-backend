package main

import (
	"context"
	"unified-go-backend/config"
	"unified-go-backend/database"
	"unified-go-backend/routes"
	"unified-go-backend/utils"
	"unified-go-backend/worker"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/sync/errgroup"
)

// @title Go REST API
// @version 1.0
// @description This is a sample server for a Go REST API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email you@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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

	// Serve Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	utils.Logger.Info("Starting server...")

	g.Go(func() error {
		return router.Run()
	})

	if err := g.Wait(); err != nil {
		utils.Logger.Fatalf("Failed to run server: %v", err)
	}
}
