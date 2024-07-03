package middleware

import (
	"context"
	"net/http"
	"unified-go-backend/database"
	"unified-go-backend/models"
	"unified-go-backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func AuthorizationMiddleware(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		email, exists := c.Get("email")
		if !exists {
			utils.Logger.Errorf("AuthorizationMiddleware: Failed to get email from context")
			c.JSON(http.StatusUnauthorized, utils.CreateErrorResponse("Unauthorized", nil))
			c.Abort()
			return
		}

		collection := database.MongoClient.Database("mdmdb").Collection("users")
		var user models.User
		err := collection.FindOne(c.Request.Context(), bson.M{"email": email}).Decode(&user)
		if err != nil {
			utils.Logger.Errorf("AuthorizationMiddleware: Error fetching user: %v", err)
			c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error fetching user", nil))
			c.Abort()
			return
		}

		if !hasRequiredPermissions(user, requiredPermissions) {
			utils.Logger.Warnf("AuthorizationMiddleware: User %s does not have required permissions", email)
			c.JSON(http.StatusForbidden, utils.CreateErrorResponse("Forbidden", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

func hasRequiredPermissions(user models.User, requiredPermissions []string) bool {
	userPermissions := make(map[string]bool)

	// Fetch roles and permissions from the user's roles
	for _, role := range user.Roles {
		var roleDoc models.Role
		err := database.MongoClient.Database("mdmdb").Collection("roles").FindOne(context.TODO(), bson.M{"name": role}).Decode(&roleDoc)
		if err == nil {
			for _, perm := range roleDoc.Permissions {
				userPermissions[perm] = true
			}
		}
	}

	// Fetch permissions from the user's access group
	if user.AccessGroup != "" {
		var accessGroup models.AccessGroup
		err := database.MongoClient.Database("mdmdb").Collection("access_groups").FindOne(context.TODO(), bson.M{"name": user.AccessGroup}).Decode(&accessGroup)
		if err == nil {
			for _, perm := range accessGroup.Permissions {
				userPermissions[perm] = true
			}
		}
	}

	for _, requiredPerm := range requiredPermissions {
		if !userPermissions[requiredPerm] {
			return false
		}
	}
	return true
}
