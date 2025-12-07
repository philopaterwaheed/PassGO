package backend

import (
	"net/http"
	"github.com/philopaterwaheed/passGO/internal/backend/config"
	"github.com/gin-gonic/gin"
)

// Run starts the Gin HTTP server
func Run() {
	router := SetupRouter()
	router.Run(":" + config.Port)
}

// SetupRouter configures and returns the Gin router
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// API routes group
	api := router.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	return router
}
