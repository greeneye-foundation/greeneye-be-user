package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/greeneye-foundation/greeneye-be-user/internal/config"
	"github.com/greeneye-foundation/greeneye-be-user/internal/models"
	"github.com/greeneye-foundation/greeneye-be-user/internal/pkg/utils"
)

type AuthService struct {
	userService *UserService
	cfg         *config.Config
	redisClient *redis.Client
	validate    *validator.Validate
}

const passwordResetKeyFormat = "password_reset:%s"

func NewAuthService(userService *UserService, cfg *config.Config, redisClient *redis.Client) *AuthService {
	return &AuthService{
		userService: userService,
		cfg:         cfg,
		redisClient: redisClient,
		validate:    validator.New(),
	}
}

// RegisterUser handles user registration logic
func (a *AuthService) RegisterUser(ctx context.Context, reg *models.UserRegistration) (*models.User, error) {
	// Check if user already exists
	existingUser, err := a.userService.GetUserByMobileNumber(ctx, reg.MobileNumber)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Create user object
	user := &models.User{
		MobileNumber: reg.MobileNumber,
		CountryCode:  reg.CountryCode,
		PasswordHash: reg.Password,
		IsVerified:   false,
		Roles:        []string{"user"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create user via UserService
	if err := a.userService.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// TODO: Send OTP for mobile verification

	return user, nil
}

// LoginUser handles user authentication and token generation
func (a *AuthService) LoginUser(ctx context.Context, login *models.UserLogin) (string, error) {
	user, err := a.userService.AuthenticateUser(ctx, login)
	if err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := a.generateJWT(user.ID.Hex())
	if err != nil {
		return "", err
	}

	// Update last login time
	user.LastLoginAt = time.Now()
	a.userService.UpdateUser(ctx, user)

	return token, nil
}

// PasswordRecovery initiates the password recovery process
func (a *AuthService) PasswordRecovery(ctx context.Context, req *models.PasswordRecoveryRequest) error {
	// Find user by mobile number
	user, err := a.userService.GetUserByMobileNumber(ctx, req.MobileNumber)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	// Generate a password reset token
	token := utils.GenerateRandomToken(32)

	// Store token in Redis with expiration (e.g., 15 minutes)
	err = a.redisClient.Set(ctx, fmt.Sprintf(passwordResetKeyFormat, token), user.ID.Hex(), 15*time.Minute).Err()
	if err != nil {
		return errors.New("failed to generate password reset token")
	}

	resetLink := fmt.Sprintf("https://yourdomain.com/reset-password?token=%s", token)
	smsMessage := fmt.Sprintf("Your password reset token is: %s", resetLink)

	if err := utils.SendSMS(user.MobileNumber, smsMessage); err != nil {
		return errors.New("failed to send SMS")
	}

	return nil
}

// ResetPassword completes the password reset process
func (a *AuthService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	// Retrieve user ID from Redis using the token
	userID, err := a.redisClient.Get(ctx, fmt.Sprintf(passwordResetKeyFormat, req.Token)).Result()
	if err == redis.Nil {
		return errors.New("invalid or expired token")
	} else if err != nil {
		return errors.New("failed to validate token")
	}

	// Get user by ID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}
	user, err := a.userService.GetUserByID(ctx, userObjectID)
	if err != nil {
		return errors.New("user not found")
	}

	// Update password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	if err := a.userService.UpdateUser(ctx, user); err != nil {
		return errors.New("failed to update password")
	}

	// Delete the reset token
	a.redisClient.Del(ctx, fmt.Sprintf(passwordResetKeyFormat, req.Token))

	return nil
}

func (a *AuthService) generateJWT(userID string) (string, error) {
	// Define token claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	return token.SignedString([]byte(a.cfg.JWT.Secret))
}

func (a *AuthService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return a.userService.GetUserByID(ctx, id)
}
