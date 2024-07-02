package controllers

import (
	"context"
	"net/http"
	"strconv"
	"unified-go-backend/config"
	"unified-go-backend/database"
	"unified-go-backend/models"
	"unified-go-backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/v1/user/profile [get]
// @Security BearerAuth
func (u *UserController) Profile(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		utils.Logger.Errorf("Profile: Failed to get email from context")
		c.JSON(http.StatusUnauthorized, utils.CreateErrorResponse("Unauthorized"))
		return
	}

	collection := database.MongoClient.Database("mdmdb").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		utils.Logger.Errorf("Profile: Error fetching user profile for email: %s, error: %v", email, err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error fetching user profile"))
		return
	}

	utils.Logger.Infof("Fetched user profile for email: %s", email)
	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.User true "User profile data"
// @Success 200 {object} map[string]string "User profile updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Error updating user profile"
// @Router /api/v1/user/profile [put]
// @Security BearerAuth
func (u *UserController) UpdateProfile(c *gin.Context) {
	email := c.MustGet("email").(string)

	var userUpdate models.User
	if err := c.BindJSON(&userUpdate); err != nil {
		utils.Logger.Errorf("UpdateProfile: Invalid request for email: %s, error: %v", email, err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request"))
		return
	}

	collection := database.MongoClient.Database("mdmdb").Collection("users")
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

// ListUsers godoc
// @Summary List all users
// @Description List all users with pagination
// @Tags user
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/v1/users [get]
// @Security BearerAuth
func (u *UserController) ListUsers(c *gin.Context) {
	collection := database.MongoClient.Database("mdmdb").Collection("users")

	// Get pagination parameters from query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	// Find users with pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	var users []models.User
	cursor, err := collection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		utils.Logger.Errorf("ListUsers: Error fetching users: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error fetching users"))
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			utils.Logger.Errorf("ListUsers: Error decoding user: %v", err)
			c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error decoding user"))
			return
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		utils.Logger.Errorf("ListUsers: Cursor error: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Cursor error"))
		return
	}

	// Get total count of users
	totalCount, err := collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		utils.Logger.Errorf("ListUsers: Error counting users: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error counting users"))
		return
	}

	utils.Logger.Infof("Fetched %d users", len(users))
	c.JSON(http.StatusOK, utils.CreatePaginatedResponse(users, page, limit, int(totalCount)))
}
