package routes

import (
	"unified-go-backend/config"
	"unified-go-backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, cfg *config.Config) {
	authController := controllers.NewAuthController(cfg)
	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
}
