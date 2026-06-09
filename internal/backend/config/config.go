package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Port        string
	Environment string

	JWTSecret     string
	JWTExpiration int

	MongoURI      string
	MongoDatabase string

	SupabaseURL    string
	SupabaseAPIKey string
	PublicBaseURL  string

	RateLimitRedisURL      string
	RateLimitRequests      int
	RateLimitWindowSeconds int

	MaxVaultsPerUser  int
	MaxVaultDataBytes int
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using default values")
	} else {
		log.Println(".env file loaded successfully")
	}
	Port = getEnv("PORT", "8080")
	Environment = getEnv("ENVIRONMENT", "development")
	JWTSecret = getEnv("JWT_SECRET", "")
	JWTExpiration = getEnvAsInt("JWT_EXPIRATION_HOURS", 24)
	MongoURI = getEnv("MONGO_URI", "mongodb://localhost:27017")
	MongoDatabase = getEnv("MONGO_DATABASE", "passgo")
	SupabaseURL = getEnv("SUPABASE_URL", "")
	SupabaseAPIKey = getEnv("SUPABASE_API_KEY", "")
	PublicBaseURL = getEnv("PASSGO_PUBLIC_BASE_URL", "")
	RateLimitRedisURL = getEnv("UPSTASH_REDIS_URL", getEnv("REDIS_URL", ""))
	RateLimitRequests = getEnvAsInt("RATE_LIMIT_REQUESTS", 20)
	RateLimitWindowSeconds = getEnvAsInt("RATE_LIMIT_WINDOW_SECONDS", 60)
	MaxVaultsPerUser = getEnvAsInt("MAX_VAULTS_PER_USER", 50)
	MaxVaultDataBytes = getEnvAsInt("MAX_VAULT_DATA_BYTES", 16*1024)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
