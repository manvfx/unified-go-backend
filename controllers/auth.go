package controllers

import (
	"context"
	"net/http"
	"time"
	"unified-go-backend/config"
	"unified-go-backend/database"
	"unified-go-backend/models"
	"unified-go-backend/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	config *config.Config
}

func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		config: cfg,
	}
}

func (a *AuthController) Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		utils.Logger.Errorf("Register: Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Errorf("Register: Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error hashing password"))
		return
	}
	user.Password = string(hashedPassword)
	user.Verified = false

	collection := database.MongoClient.Database("testdb").Collection("users")
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		utils.Logger.Errorf("Register: Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error creating user"))
		return
	}

	// Generate and save verification code
	verificationCode := utils.GenerateVerificationCode()
	err = database.RedisClient.Set(context.TODO(), user.Email, verificationCode, 10*time.Minute).Err()
	if err != nil {
		utils.Logger.Errorf("Register: Error saving verification code: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error saving verification code"))
		return
	}

	// Add email to the verification queue
	err = database.RedisClient.LPush(context.TODO(), "email_verification_queue", user.Email).Err()
	if err != nil {
		utils.Logger.Errorf("Register: Error adding email to verification queue: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error adding email to verification queue"))
		return
	}

	utils.Logger.Infof("User registered successfully: %s", user.Email)
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully. Please check your email for the verification code."})
}

func (a *AuthController) VerifyEmail(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := c.BindJSON(&request); err != nil {
		utils.Logger.Errorf("VerifyEmail: Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request"))
		return
	}

	storedCode, err := database.RedisClient.Get(context.TODO(), request.Email).Result()
	if err != nil || storedCode != request.Code {
		utils.Logger.Errorf("VerifyEmail: Invalid or expired verification code for email: %s", request.Email)
		c.JSON(http.StatusUnauthorized, utils.CreateErrorResponse("Invalid or expired verification code"))
		return
	}

	collection := database.MongoClient.Database("testdb").Collection("users")
	filter := bson.M{"email": request.Email}
	update := bson.M{
		"$set": bson.M{
			"verified":    true,
			"verified_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		utils.Logger.Errorf("VerifyEmail: Error updating user verification status: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error updating user verification status"))
		return
	}

	utils.Logger.Infof("Email verified successfully: %s", request.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (a *AuthController) Login(c *gin.Context) {
	var reqUser models.User
	if err := c.BindJSON(&reqUser); err != nil {
		utils.Logger.Errorf("Login: Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request"))
		return
	}

	collection := database.MongoClient.Database("testdb").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": reqUser.Email}).Decode(&user)
	if err != nil {
		utils.Logger.Errorf("Login: Invalid email or password: %s", reqUser.Email)
		c.JSON(http.StatusUnauthorized, utils.CreateErrorResponse("Invalid email or password"))
		return
	}

	if !user.Verified {
		utils.Logger.Warnf("Login: Email not verified: %s", reqUser.Email)
		c.JSON(http.StatusUnauthorized, utils.CreateErrorResponse("Email not verified"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password))
	if err != nil {
		utils.Logger.Errorf("Login: Invalid email or password: %s", reqUser.Email)
		c.JSON(http.StatusUnauthorized, utils.CreateErrorResponse("Invalid email or password"))
		return
	}

	// Capture system information
	userIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	user.LastLogin = time.Now()
	user.LastLoginIP = userIP
	user.LastLoginAgent = userAgent

	// Update user login info in the database
	filter := bson.M{"email": user.Email}
	update := bson.M{
		"$set": bson.M{
			"last_login":       user.LastLogin,
			"last_login_ip":    user.LastLoginIP,
			"last_login_agent": user.LastLoginAgent,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		utils.Logger.Errorf("Login: Error updating user login info: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error updating user login info"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.config.JwtSecret))
	if err != nil {
		utils.Logger.Errorf("Login: Error creating token: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error creating token"))
		return
	}

	utils.Logger.Infof("User logged in successfully: %s", user.Email)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
