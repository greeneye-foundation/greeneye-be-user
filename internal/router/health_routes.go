package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (r *Router) setupHealthRoutes() {
	health := r.router.Group("/health")

	health.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	health.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}