package routes

import (
	"unified-go-backend/config"
	"unified-go-backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, cfg *config.Config) {
	authController := controllers.NewAuthController(cfg)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", authController.Register)
		v1.POST("/login", authController.Login)
		v1.POST("/verify-email", authController.VerifyEmail)
	}
}
