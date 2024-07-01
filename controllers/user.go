package controllers

import (
	"context"
	"net/http"
	"unified-go-backend/config"
	"unified-go-backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserController struct {
	config      *config.Config
	mongoClient *mongo.Client
}

func NewUserController(cfg *config.Config) *UserController {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		panic(err)
	}

	return &UserController{
		config:      cfg,
		mongoClient: client,
	}
}

func (u *UserController) Profile(c *gin.Context) {
	username := c.MustGet("username").(string)

	collection := u.mongoClient.Database("testdb").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
