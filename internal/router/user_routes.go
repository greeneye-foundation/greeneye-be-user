package router

import (
	"github.com/gin-gonic/gin"

	"github.com/greeneye-foundation/greeneye-be-user/internal/handlers"
	"github.com/greeneye-foundation/greeneye-be-user/internal/services"
)

func userRoutes(r *gin.RouterGroup, userService *services.UserService) {
	userHandler := handlers.NewUserHandler(userService)

	userGroup := r.Group("/user")
	{
		userGroup.GET("/profile/:id", userHandler.GetProfile)
		// Add more user-related routes here
	}
}
