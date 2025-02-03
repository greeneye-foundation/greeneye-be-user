package logger

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// InitLogger initializes a singleton zap logger
func InitLogger(environment string) *zap.Logger {
	once.Do(func() {
		var err error
		var config zap.Config

		// Different configurations for development and production
		if environment == "production" {
			config = zap.NewProductionConfig()
			config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			config.OutputPaths = []string{"./logs/app.log"}
			config.ErrorOutputPaths = []string{"./logs/errors.log"}
		} else {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		// Custom sampling for high-frequency logs
		config.Sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}

		logger, err = config.Build(zap.AddCallerSkip(1))
		if err != nil {
			panic(err)
		}
	})

	return logger
}

// GetLogger returns the initialized logger
func GetLogger() *zap.Logger {
	if logger == nil {
		return InitLogger("development")
	}
	return logger
}

// WithError creates a logger with an error field
func WithError(err error) *zap.Logger {
	return GetLogger().With(zap.Error(err))
}

// LoggerMiddleware creates a gin middleware for logging requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Log details after request
		duration := time.Since(start)
		logger := GetLogger().With(
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", duration),
			zap.String("client_ip", c.ClientIP()),
		)

		if len(c.Errors) > 0 {
			logger.Error("Request failed", zap.Any("errors", c.Errors))
		} else {
			logger.Info("Request processed")
		}
	}
}
