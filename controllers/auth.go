package controllers

import (
	"context"
	"net/http"
	"time"
	"unified-go-backend/config"
	"unified-go-backend/models"
	"unified-go-backend/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	config      *config.Config
	mongoClient *mongo.Client
}

func NewAuthController(cfg *config.Config) *AuthController {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		panic(err)
	}

	return &AuthController{
		config:      cfg,
		mongoClient: client,
	}
}

func (a *AuthController) Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error hashing password"))
		return
	}
	user.Password = string(hashedPassword)

	collection := a.mongoClient.Database("testdb").Collection("users")
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error creating user"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (a *AuthController) Login(c *gin.Context) {
	var reqUser models.User
	if err := c.BindJSON(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request"))
		return
	}

	collection := a.mongoClient.Database("testdb").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": reqUser.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.CreateErrorResponse("Invalid username or password"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.CreateErrorResponse("Invalid username or password"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.config.JwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error creating token"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
