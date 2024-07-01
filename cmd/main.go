package main

import (
	"unified-go-backend/config"
	"unified-go-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	router := gin.Default()

	routes.AuthRoutes(router, cfg)
	routes.UserRoutes(router, cfg)

	router.Run()
}
