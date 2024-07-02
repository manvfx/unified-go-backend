package routes

import (
	"unified-go-backend/config"
	"unified-go-backend/controllers"
	"unified-go-backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, cfg *config.Config) {
	userController := controllers.NewUserController(cfg)

	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(cfg))
	{
		v1.GET("/user/profile", userController.Profile)
		v1.PUT("/user/profile", userController.UpdateProfile)
		v1.PUT("/user/:id", userController.UpdateUser)
		v1.DELETE("/user/:id", userController.DeleteUser)
		v1.GET("/users", userController.ListUsers)
	}
}
