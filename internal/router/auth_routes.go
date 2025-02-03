package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/greeneye-foundation/greeneye-be-user/internal/config"
	"github.com/greeneye-foundation/greeneye-be-user/internal/handlers"
	"github.com/greeneye-foundation/greeneye-be-user/internal/middleware"
	"github.com/greeneye-foundation/greeneye-be-user/internal/services"
)

func authRoutes(r *gin.RouterGroup, userService *services.UserService, cfg *config.Config, redisClient *redis.Client) {
	authService := services.NewAuthService(userService, cfg, redisClient)
	authHandler := handlers.NewAuthHandler(authService, cfg, redisClient)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/password-recovery", authHandler.PasswordRecovery)
		authGroup.POST("/reset-password", authHandler.ResetPassword)
	}

	// Protected routes example
	protected := r.Group("/protected")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		protected.GET("/profile", authHandler.Profile)
	}
}
