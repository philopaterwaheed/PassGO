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
