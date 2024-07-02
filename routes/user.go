package routes

import (
	"unified-go-backend/config"
	"unified-go-backend/controllers"
	"unified-go-backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, cfg *config.Config) {
	userController := controllers.NewUserController(cfg)

	auth := router.Group("/user")
	auth.Use(middleware.AuthMiddleware(cfg))
	auth.GET("/profile", userController.Profile)
	auth.PUT("/profile", userController.UpdateProfile)
}
