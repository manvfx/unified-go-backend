package routes

import (
	"unified-go-backend/config"
	"unified-go-backend/controllers"
	"unified-go-backend/middleware"

	"github.com/gin-gonic/gin"
)

func AccessGroupRoutes(router *gin.Engine, cfg *config.Config) {
	accessGroupController := controllers.NewAccessGroupController()

	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(cfg))
	{
		v1.POST("/access_groups", middleware.AuthorizationMiddleware("create_access_group"), accessGroupController.CreateAccessGroup)
		v1.GET("/access_groups", middleware.AuthorizationMiddleware("list_access_groups"), accessGroupController.ListAccessGroups)
		v1.PUT("/access_groups/:id", middleware.AuthorizationMiddleware("update_access_group"), accessGroupController.UpdateAccessGroup)
		v1.DELETE("/access_groups/:id", middleware.AuthorizationMiddleware("delete_access_group"), accessGroupController.DeleteAccessGroup)
	}
}
