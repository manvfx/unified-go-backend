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
		v1.PUT("/user/:id", middleware.AuthorizationMiddleware("update_user"), userController.UpdateUser)
		v1.DELETE("/user/:id", middleware.AuthorizationMiddleware("delete_user"), userController.DeleteUser)
		v1.GET("/users", middleware.AuthorizationMiddleware("list_users"), userController.ListUsers)
	}
}
