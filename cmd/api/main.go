package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/greeneye-foundation/greeneye-be-user/internal/config"
	"github.com/greeneye-foundation/greeneye-be-user/internal/pkg/logger"
	"github.com/greeneye-foundation/greeneye-be-user/internal/router"
)

// cmd/api/main.go
// @title           GreenEye User Management API
// @version         1.0
// @description     User management microservice
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http https
func main() {
	// Initialize logger
	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "development"
	}
	log := logger.InitLogger(environment)
	defer func() {
		_ = log.Sync()
	}()

	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Error loading config", zap.Error(err))
	}

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize MongoDB
	mongoClient, err := config.InitMongoDB(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Error("Error disconnecting MongoDB", zap.Error(err))
		}
	}()

	// Initialize Redis
	redisClient, err := config.InitRedis(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Setup router
	r := router.NewRouter(cfg, mongoClient, redisClient)

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r.Engine(),
	}

	// Start server in a goroutine
	go func() {
		log.Info("Starting server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")
}
