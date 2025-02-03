package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/greeneye-foundation/greeneye-be-user/internal/config"
	"github.com/greeneye-foundation/greeneye-be-user/internal/handlers"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Router struct {
	router *gin.Engine
	config *config.Config
	db     *mongo.Client
	redis  *redis.Client
}

func NewRouter(
	cfg *config.Config,
	db *mongo.Client,
	redisClient *redis.Client,
) *Router {
	// Set Gin mode based on environment
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Create Router struct
	r := &Router{
		router: router,
		config: cfg,
		db:     db,
		redis:  redisClient,
	}

	// Setup routes
	r.setupMiddleware()
	r.setupHealthRoutes()
	r.setupRoutes()

	return r
}

func (r *Router) setupRoutes() {
	// Setup main application routes
	handlers.SetupRoutes(r.router, r.db, r.config, r.redis)
}

func (r *Router) setupMiddleware() {
	// Add any global middleware
	// r.router.Use(someMiddleware())
}

func (r *Router) Start() error {
	// Start the server
	return r.router.Run(fmt.Sprintf(":%s", r.config.Server.Port))
}

// Optional: Method to get the underlying Gin engine
func (r *Router) Engine() *gin.Engine {
	return r.router
}
