package controllers

import (
	"context"
	"net/http"
	"unified-go-backend/database"
	"unified-go-backend/models"
	"unified-go-backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessGroupController struct{}

func NewAccessGroupController() *AccessGroupController {
	return &AccessGroupController{}
}

// CreateAccessGroup godoc
// @Summary Create a new access group
// @Description Create a new access group with roles and permissions
// @Tags access_group
// @Accept json
// @Produce json
// @Param access_group body models.AccessGroup true "Access group data"
// @Success 201 {object} models.AccessGroup
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/access_groups [post]
// @Security BearerAuth
func (a *AccessGroupController) CreateAccessGroup(c *gin.Context) {
	var accessGroup models.AccessGroup
	if err := c.BindJSON(&accessGroup); err != nil {
		utils.Logger.Errorf("CreateAccessGroup: Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request", nil))
		return
	}

	// Validate the access group request
	if err := utils.ValidateStruct(accessGroup); err != nil {
		validationErrors := utils.FormatValidationError(err)
		utils.Logger.Errorf("CreateAccessGroup: Validation error: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Validation error", validationErrors))
		return
	}

	collection := database.MongoClient.Database("mdmdb").Collection("access_groups")
	result, err := collection.InsertOne(context.TODO(), accessGroup)
	if err != nil {
		utils.Logger.Errorf("CreateAccessGroup: Error creating access group: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error creating access group", nil))
		return
	}

	accessGroup.ID = result.InsertedID.(primitive.ObjectID)
	utils.Logger.Infof("Access group created successfully: %s", accessGroup.Name)
	c.JSON(http.StatusCreated, accessGroup)
}

// ListAccessGroups godoc
// @Summary List all access groups
// @Description List all access groups
// @Tags access_group
// @Produce json
// @Success 200 {array} models.AccessGroup
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/access_groups [get]
// @Security BearerAuth
func (a *AccessGroupController) ListAccessGroups(c *gin.Context) {
	collection := database.MongoClient.Database("mdmdb").Collection("access_groups")

	var accessGroups []models.AccessGroup
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		utils.Logger.Errorf("ListAccessGroups: Error fetching access groups: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error fetching access groups", nil))
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var accessGroup models.AccessGroup
		if err := cursor.Decode(&accessGroup); err != nil {
			utils.Logger.Errorf("ListAccessGroups: Error decoding access group: %v", err)
			c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error decoding access group", nil))
			return
		}
		accessGroups = append(accessGroups, accessGroup)
	}

	if err := cursor.Err(); err != nil {
		utils.Logger.Errorf("ListAccessGroups: Cursor error: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Cursor error", nil))
		return
	}

	utils.Logger.Infof("Fetched %d access groups", len(accessGroups))
	c.JSON(http.StatusOK, accessGroups)
}

// UpdateAccessGroup godoc
// @Summary Update an access group
// @Description Update an access group's details
// @Tags access_group
// @Accept json
// @Produce json
// @Param id path string true "Access Group ID"
// @Param access_group body models.AccessGroup true "Access group data"
// @Success 200 {object} map[string]string "message": "Access group updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 404 {object} utils.ErrorResponse "Access group not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/access_groups/{id} [put]
// @Security BearerAuth
func (a *AccessGroupController) UpdateAccessGroup(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Logger.Errorf("UpdateAccessGroup: Invalid access group ID: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid access group ID", nil))
		return
	}

	var accessGroup models.AccessGroup
	if err := c.BindJSON(&accessGroup); err != nil {
		utils.Logger.Errorf("UpdateAccessGroup: Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request", nil))
		return
	}

	// Validate the access group request
	if err := utils.ValidateStruct(accessGroup); err != nil {
		validationErrors := utils.FormatValidationError(err)
		utils.Logger.Errorf("UpdateAccessGroup: Validation error: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Validation error", validationErrors))
		return
	}

	collection := database.MongoClient.Database("mdmdb").Collection("access_groups")
	update := bson.M{
		"$set": bson.M{
			"name":        accessGroup.Name,
			"roles":       accessGroup.Roles,
			"permissions": accessGroup.Permissions,
		},
	}

	filter := bson.M{"_id": objectId}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		utils.Logger.Errorf("UpdateAccessGroup: Error updating access group: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error updating access group", nil))
		return
	}
	if result.MatchedCount == 0 {
		utils.Logger.Errorf("UpdateAccessGroup: Access group not found with ID: %s", id)
		c.JSON(http.StatusNotFound, utils.CreateErrorResponse("Access group not found", nil))
		return
	}

	utils.Logger.Infof("Access group updated successfully: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Access group updated successfully"})
}

// DeleteAccessGroup godoc
// @Summary Delete an access group
// @Description Delete an access group by ID
// @Tags access_group
// @Produce json
// @Param id path string true "Access Group ID"
// @Success 200 {object} map[string]string "message": "Access group deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 404 {object} utils.ErrorResponse "Access group not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/access_groups/{id} [delete]
// @Security BearerAuth
func (a *AccessGroupController) DeleteAccessGroup(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Logger.Errorf("DeleteAccessGroup: Invalid access group ID: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid access group ID", nil))
		return
	}

	collection := database.MongoClient.Database("mdmdb").Collection("access_groups")
	filter := bson.M{"_id": objectId}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		utils.Logger.Errorf("DeleteAccessGroup: Error deleting access group: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error deleting access group", nil))
		return
	}
	if result.DeletedCount == 0 {
		utils.Logger.Errorf("DeleteAccessGroup: Access group not found with ID: %s", id)
		c.JSON(http.StatusNotFound, utils.CreateErrorResponse("Access group not found", nil))
		return
	}

	utils.Logger.Infof("Access group deleted successfully: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Access group deleted successfully"})
}
