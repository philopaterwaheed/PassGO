package backend

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/philopaterwaheed/passGO/internal/backend/config"
	"github.com/philopaterwaheed/passGO/internal/backend/database"
)

// Run starts the Gin HTTP server
func Run() {
	// Initialize MongoDB connection
	if err := database.Connect(config.MongoURI, config.MongoDatabase); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.Disconnect()

	router := SetupRouter()
	router.Run(":" + config.Port)
}

// SetupRouter configures and returns the Gin router
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		status := "healthy"
		dbStatus := "connected"

		// Check database connection
		if err := database.HealthCheck(); err != nil {
			status = "degraded"
			dbStatus = "disconnected"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   status,
			"database": dbStatus,
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
