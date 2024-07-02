package controllers

import (
	"context"
	"net/http"
	"unified-go-backend/config"
	"unified-go-backend/database"
	"unified-go-backend/models"
	"unified-go-backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// UserController handles user-related operations.
type UserController struct {
	config *config.Config
}

// NewUserController creates a new UserController.
func NewUserController(cfg *config.Config) *UserController {
	return &UserController{
		config: cfg,
	}
}

// Profile godoc
// @Summary Get user profile
// @Description Get the authenticated user's profile
// @Tags user
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} utils.ErrorResponse
// @Router /user/profile [get]
// @Security BearerAuth
func (u *UserController) Profile(c *gin.Context) {
	email := c.MustGet("email").(string)

	collection := database.MongoClient.Database("testdb").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		utils.Logger.Errorf("Profile: Error fetching user profile for email: %s, error: %v", email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user profile"})
		return
	}

	utils.Logger.Infof("Fetched user profile for email: %s", email)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.User true "User profile data"
// @Success 200 {object} gin.H{"message": string}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /user/profile [put]
// @Security BearerAuth
func (u *UserController) UpdateProfile(c *gin.Context) {
	email := c.MustGet("email").(string)

	var userUpdate models.User
	if err := c.BindJSON(&userUpdate); err != nil {
		utils.Logger.Errorf("UpdateProfile: Invalid request for email: %s, error: %v", email, err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request"))
		return
	}

	collection := database.MongoClient.Database("testdb").Collection("users")
	update := bson.M{
		"$set": bson.M{
			"username": userUpdate.Username,
			"password": userUpdate.Password, // This assumes the password is already hashed
		},
	}

	filter := bson.M{"email": email}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		utils.Logger.Errorf("UpdateProfile: Error updating user profile for email: %s, error: %v", email, err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error updating user profile"))
		return
	}

	utils.Logger.Infof("User profile updated successfully for email: %s", email)
	c.JSON(http.StatusOK, gin.H{"message": "User profile updated successfully"})
}
