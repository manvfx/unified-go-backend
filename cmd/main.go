package main

import (
	"context"
	// "flag"
	// "fmt"
	"unified-go-backend/config"
	"unified-go-backend/database"
	_ "unified-go-backend/docs"
	"unified-go-backend/middleware"
	"unified-go-backend/routes"
	"unified-go-backend/utils"
	"unified-go-backend/worker"

	// "unified-go-backend/seed"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/sync/errgroup"
)

// @title Unified Go Backend
// @version 1.0
// @description This is a sample server for a Unified Go Backend.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email manvfx@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// seedFlag := flag.Bool("seed", false, "Seed the database with dummy data")
	// flag.Parse()

	cfg := config.LoadConfig()

	// Initialize logger with Elasticsearch URL
	utils.InitLogger()

	// Connect to the database
	database.ConnectDB(cfg)
	defer database.DisconnectDB()

	// Connect to Redis
	database.ConnectRedis(cfg)
	defer database.DisconnectRedis()

	// if *seedFlag {
	//     seed.SeedData(cfg)
	//     fmt.Println("Database seeded with dummy data.")
	//     return
	// }

	g, ctx := errgroup.WithContext(context.Background())

	// Start the email verification job worker
	g.Go(func() error {
		return worker.ProcessEmailVerificationJobs(ctx, cfg)
	})

	router := gin.Default()

	// Create a new RateLimiter instance and apply the rate limiter middleware globally
	rateLimiter := middleware.NewRateLimiter()
	router.Use(middleware.RateLimiterMiddleware(rateLimiter))

	// Setup routes with versioning
	routes.AuthRoutes(router, cfg)
	routes.UserRoutes(router, cfg)
	routes.AccessGroupRoutes(router, cfg)

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
