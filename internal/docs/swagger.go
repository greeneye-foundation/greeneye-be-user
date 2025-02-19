// internal/docs/swagger.go
package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/greeneye-foundation/greeneye-be-user/docs" //Import the docs generated file
)

// @title           GreenEye User Management API
// @version         1.0
// @description     User management and authentication API
// @host            localhost:8080
// @BasePath        /api/v1

func SetupSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
