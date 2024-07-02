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
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// AuthController handles authentication-related operations.
type AuthController struct {
	config *config.Config
}

// NewAuthController creates a new AuthController.
func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		config: cfg,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User to register"
// @Success 201 {object} map[string]string "message": "User created successfully. Please check your email for the verification code."
// @Failure 400 {object} utils.ErrorResponse "Validation error"
// @Failure 409 {object} utils.ErrorResponse "User already exists"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/register [post]
func (a *AuthController) Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		utils.Logger.Errorf("Register: Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request", nil))
		return
	}

	// Validate the user request
	if err := utils.ValidateStruct(user); err != nil {
		validationErrors := utils.FormatValidationError(err)
		utils.Logger.Errorf("Register: Validation error: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Validation error", validationErrors))
		return
	}

	collection := database.MongoClient.Database("mdmdb").Collection("users")

	// Check for duplicate user
	var existingUser models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		utils.Logger.Errorf("Register: User already exists with email: %s", user.Email)
		c.JSON(http.StatusConflict, utils.CreateErrorResponse("User already exists", nil))
		return
	}
	if err != mongo.ErrNoDocuments {
		utils.Logger.Errorf("Register: Error checking for duplicate user: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error checking for duplicate user", nil))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Errorf("Register: Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error hashing password", nil))
		return
	}
	user.Password = string(hashedPassword)
	user.Verified = false

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		utils.Logger.Errorf("Register: Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error creating user", nil))
		return
	}

	// Generate and save verification code
	verificationCode := utils.GenerateVerificationCode()
	err = database.RedisClient.Set(context.TODO(), user.Email, verificationCode, 10*time.Minute).Err()
	if err != nil {
		utils.Logger.Errorf("Register: Error saving verification code: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error saving verification code", nil))
		return
	}

	// Add email to the verification queue
	err = database.RedisClient.LPush(context.TODO(), "email_verification_queue", user.Email).Err()
	if err != nil {
		utils.Logger.Errorf("Register: Error adding email to verification queue: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error adding email to verification queue", nil))
		return
	}

	utils.Logger.Infof("User registered successfully: %s", user.Email)
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully. Please check your email for the verification code."})
}

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify a user's email address with a verification code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.VerifyEmailRequest true "Email verification request"
// @Success 200 {object} map[string]string "message": "Email verified successfully"
// @Failure 400 {object} utils.ErrorResponse "Validation error"
// @Failure 401 {object} utils.ErrorResponse "Invalid or expired verification code"
// @Failure 500 {object} utils.ErrorResponse "Error updating user verification status"
// @Router /api/v1/verify-email [post]
func (a *AuthController) VerifyEmail(c *gin.Context) {
	var request models.VerifyEmailRequest
	if err := c.BindJSON(&request); err != nil {
		utils.Logger.Errorf("VerifyEmail: Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request", nil))
		return
	}

	// Validate the request
	if err := utils.ValidateStruct(request); err != nil {
		validationErrors := utils.FormatValidationError(err)
		utils.Logger.Errorf("VerifyEmail: Validation error: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Validation error", validationErrors))
		return
	}

	storedCode, err := database.RedisClient.Get(context.TODO(), request.Email).Result()
	if err != nil || storedCode != request.Code {
		utils.Logger.Errorf("VerifyEmail: Invalid or expired verification code for email: %s", request.Email)
		c.JSON(http.StatusUnauthorized, utils.CreateErrorResponse("Invalid or expired verification code", nil))
		return
	}

	collection := database.MongoClient.Database("mdmdb").Collection("users")
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
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error updating user verification status", nil))
		return
	}

	utils.Logger.Infof("Email verified successfully: %s", request.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "User credentials"
// @Success 200 {object} models.LoginResponse "Returns a token on successful login"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Invalid email or password"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/login [post]
func (a *AuthController) Login(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.BindJSON(&loginRequest); err != nil {
		utils.Logger.Errorf("Login: Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Invalid request", nil))
		return
	}

	// Validate the login request
	if err := utils.ValidateStruct(loginRequest); err != nil {
		validationErrors := utils.FormatValidationError(err)
		utils.Logger.Errorf("Login: Validation error: %v", err)
		c.JSON(http.StatusBadRequest, utils.CreateErrorResponse("Validation error", validationErrors))
		return
	}

	collection := database.MongoClient.Database("mdmdb").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": loginRequest.Email}).Decode(&user)
	if err != nil {
		utils.Logger.Errorf("Login: Invalid email or password: %s", loginRequest.Email)
		c.JSON(http.StatusUnauthorized, utils.CreateErrorResponse("Invalid email or password", nil))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		utils.Logger.Errorf("Login: Invalid email or password: %s", loginRequest.Email)
		c.JSON(http.StatusUnauthorized, utils.CreateErrorResponse("Invalid email or password", nil))
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.config.JwtSecret))
	if err != nil {
		utils.Logger.Errorf("Login: Error creating token: %v", err)
		c.JSON(http.StatusInternalServerError, utils.CreateErrorResponse("Error creating token", nil))
		return
	}

	utils.Logger.Infof("User logged in successfully: %s", user.Email)
	c.JSON(http.StatusOK, models.LoginResponse{Token: tokenString})
}
