package backend

import (
	"context"
	"log"

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
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowCredentials = true
	corsConfig.AllowOriginFunc = func(origin string) bool {
		return true
	}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Master-Password"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	router.Use(cors.New(corsConfig))

	// API routes group
	api := router.Group("/api")
	rateLimiter, err := middleware.NewIPRateLimiter(middleware.RateLimitOptions{
		RedisURL:      config.RateLimitRedisURL,
		Requests:      config.RateLimitRequests,
		WindowSeconds: config.RateLimitWindowSeconds,
	})
	if err != nil {
		log.Printf("Warning: Rate limiter not initialized: %v", err)
	} else if rateLimiter != nil {
		api.Use(rateLimiter)
		log.Printf(
			"Rate limiter enabled: %d requests per %d seconds",
			config.RateLimitRequests,
			config.RateLimitWindowSeconds,
		)
	}
	{
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
				auth.GET("/reset-password", authHandler.ResetPasswordPage)
				auth.POST("/update-password", authHandler.UpdatePassword)
				auth.POST("/refresh", authHandler.RefreshToken)

				// Protected auth routes
				auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
				auth.PUT("/account-password", middleware.AuthMiddleware(), authHandler.UpdateAccountPassword)
			}
		}

		// Vault routes
		vaultHandler := handlers.NewVaultHandler()
		vaultsGroup := api.Group("/vaults")
		vaultsGroup.Use(middleware.AuthMiddleware())
		{
			vaultsGroup.POST("", vaultHandler.CreateVault)
			vaultsGroup.GET("", vaultHandler.GetVaults)
			vaultsGroup.PUT("/master-password", vaultHandler.UpdateMasterPassword)
			vaultsGroup.GET("/:id", vaultHandler.GetVault)
			vaultsGroup.PUT("/:id", vaultHandler.UpdateVault)
			vaultsGroup.DELETE("/:id", vaultHandler.DeleteVault)
		}
	}

	return router
}
