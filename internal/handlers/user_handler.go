package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greeneye-foundation/greeneye-be-user/internal/models"
	"github.com/greeneye-foundation/greeneye-be-user/internal/pkg/errors"
	"github.com/greeneye-foundation/greeneye-be-user/internal/pkg/logger"
	"github.com/greeneye-foundation/greeneye-be-user/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterUser godoc
// @Summary Create a new user
// @Description Register a new user in the system
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user body models.RegisterUserRequest true "User Creation Details"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} errors.ErrorResponse
// @Router /users [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	log := logger.GetLogger()
	var registration models.UserRegistration

	// Bind and validate input
	if err := c.ShouldBindJSON(&registration); err != nil {
		log.Error("Invalid registration input", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.New(
			http.StatusBadRequest,
			"Invalid input",
			err.Error(),
		))
		return
	}

	// Create user
	err := h.userService.CreateUser(c.Request.Context(), &models.User{
		MobileNumber: registration.MobileNumber,
		// Other fields...
	})

	if err != nil {
		// Check if it's a known custom error
		if customErr, ok := err.(*errors.CustomError); ok {
			c.JSON(customErr.Code, customErr)
		} else {
			// Log unknown errors
			log.Error("User creation failed", zap.Error(err))
			c.JSON(
				http.StatusInternalServerError,
				errors.ErrInternalServer,
			)
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers(c.Request.Context())
	if err != nil {
		logger.GetLogger().Error("Failed to get users", zap.Error(err))
		c.JSON(errors.GetHTTPStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), objectID)
	if err != nil {
		logger.GetLogger().Error("Failed to get user profile", zap.Error(err))
		c.JSON(errors.GetHTTPStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
