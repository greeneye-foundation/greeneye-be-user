package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/greeneye-foundation/greeneye-be-user/internal/config"
	"github.com/greeneye-foundation/greeneye-be-user/internal/models"
	"github.com/greeneye-foundation/greeneye-be-user/internal/pkg/errors"
	"github.com/greeneye-foundation/greeneye-be-user/internal/pkg/logger"
	"github.com/greeneye-foundation/greeneye-be-user/internal/pkg/utils"
	"github.com/greeneye-foundation/greeneye-be-user/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	cfg         *config.Config
	redisClient *redis.Client
}

func NewAuthHandler(authService *services.AuthService, cfg *config.Config, redisClient *redis.Client) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
		redisClient: redisClient,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var reg models.UserRegistration
	if err := c.ShouldBindJSON(&reg); err != nil {
		logger.GetLogger().Error("Invalid registration input", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.New(
			http.StatusBadRequest,
			"Invalid input",
			err.Error(),
		))
		return
	}

	// Validate input
	if err := utils.ValidateStruct(reg); err != nil {
		logger.GetLogger().Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.New(
			http.StatusBadRequest,
			"Validation error",
			err.Error(),
		))
		return
	}

	// Call service to register user
	user, err := h.authService.RegisterUser(c.Request.Context(), &reg)
	if err != nil {
		logger.GetLogger().Error("User registration failed", zap.Error(err))
		c.JSON(errors.GetHTTPStatusCode(err), err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var login models.UserLogin
	if err := c.ShouldBindJSON(&login); err != nil {
		logger.GetLogger().Error("Invalid login input", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.New(
			http.StatusBadRequest,
			"Invalid input",
			err.Error(),
		))
		return
	}

	// Validate input
	if err := utils.ValidateStruct(login); err != nil {
		logger.GetLogger().Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.New(
			http.StatusBadRequest,
			"Validation error",
			err.Error(),
		))
		return
	}

	// Authenticate user
	token, err := h.authService.LoginUser(c.Request.Context(), &login)
	if err != nil {
		logger.GetLogger().Error("User authentication failed", zap.Error(err))
		c.JSON(errors.GetHTTPStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// PasswordRecovery handles password recovery requests
func (h *AuthHandler) PasswordRecovery(c *gin.Context) {
	var req models.PasswordRecoveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error("Invalid password recovery input", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.New(
			http.StatusBadRequest,
			"Invalid input",
			err.Error(),
		))
		return
	}

	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		logger.GetLogger().Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.New(
			http.StatusBadRequest,
			"Validation error",
			err.Error(),
		))
		return
	}

	// Initiate password recovery
	err := h.authService.PasswordRecovery(c.Request.Context(), &req)
	if err != nil {
		logger.GetLogger().Error("Password recovery failed", zap.Error(err))
		c.JSON(errors.GetHTTPStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password recovery initiated. Please check your messages.",
	})
}

// ResetPassword handles password reset requests
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error("Invalid reset password input", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.New(
			http.StatusBadRequest,
			"Invalid input",
			err.Error(),
		))
		return
	}

	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		logger.GetLogger().Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.New(
			http.StatusBadRequest,
			"Validation error",
			err.Error(),
		))
		return
	}

	// Perform password reset
	err := h.authService.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		logger.GetLogger().Error("Password reset failed", zap.Error(err))
		c.JSON(errors.GetHTTPStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been reset successfully.",
	})
}

// Profile retrieves the authenticated user's profile
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), objectID)
	if err != nil {
		logger.GetLogger().Error("Fetching user profile failed", zap.Error(err))
		c.JSON(errors.GetHTTPStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
