package backend

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/philopaterwaheed/passGO/internal/backend/config"
	"github.com/philopaterwaheed/passGO/internal/backend/database"
	"github.com/philopaterwaheed/passGO/internal/backend/handlers"
	"github.com/philopaterwaheed/passGO/internal/backend/middleware"
)

// Run starts the Gin HTTP server
func Run() {
	// Initialize MongoDB connection
	if err := database.Connect(config.MongoURI, config.MongoDatabase); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.Disconnect()

	// Initialize database indexes
	ctx := context.Background()
	userRepo := database.NewUserRepository()
	if err := userRepo.CreateIndexes(ctx); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	} else {
		log.Println("Database indexes created successfully")
	}

	router := SetupRouter()
	router.Run(":" + config.Port)
}

// SetupRouter configures and returns the Gin router
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = false
	config.AllowCredentials = true
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	router.Use(cors.New(config))

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

		// Auth routes (public)
		authHandler, err := handlers.NewAuthHandler()
		if err != nil {
			log.Printf("Warning: Auth handler not initialized (Supabase not configured): %v", err)
		} else {
			auth := api.Group("/auth")
			{
				auth.POST("/signup", authHandler.Signup)
				auth.POST("/login", authHandler.Login)
				auth.GET("/verify-email", authHandler.VerifyEmail)
				auth.POST("/verify-hash", authHandler.VerifyHash)
				auth.POST("/resend-verification", authHandler.ResendVerification)
				auth.POST("/forgot-password", authHandler.ForgotPassword)
				auth.POST("/refresh", authHandler.RefreshToken)

				// Protected auth routes
				auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
			}
		}

		// User routes
		userHandler := handlers.NewUserHandler()
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("/email/:email", userHandler.GetUserByEmail)
		}
	}

	return router
}
