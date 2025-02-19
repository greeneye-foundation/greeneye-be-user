// internal/docs/swagger.go
package docs

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/greeneye-foundation/greeneye-be-user/docs/swagger" // swagger docs
)

// @title           GreenEye User Management API
// @version         1.0
// @description     User management and authentication API
// @host            localhost:8080
// @BasePath        /api/v1

func SetupSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Example Annotated Handler
// @Summary     Register User
// @Description Create a new user account
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       request body models.UserRegistration true "User Registration"
// @Success     200 {object} models.UserResponse
// @Failure     400 {object} models.ErrorResponse
// @Router      /auth/register [post]
func RegisterUserHandler(c *gin.Context) {
	// Handler implementation
}
