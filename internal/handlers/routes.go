package handlers

import (
	"github.com/greeneye-foundation/greeneye-be-user/internal/middleware"
	"github.com/greeneye-foundation/greeneye-be-user/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/greeneye-foundation/greeneye-be-user/internal/config"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(router *gin.Engine, db *mongo.Client, cfg *config.Config, redis *redis.Client) {
	// Initialize services
	userService := services.NewUserService(db, cfg.MongoDB.Database)
	authService := services.NewAuthService(userService, cfg, redis)

	// Initialize handlers
	authHandler := NewAuthHandler(authService, cfg, redis)
	userHandler := NewUserHandler(userService)

	// API group
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// User routes
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			users.GET("/", middleware.CacheMiddleware(redis), userHandler.GetUsers)
			// Add other user routes
		}
	}
}
