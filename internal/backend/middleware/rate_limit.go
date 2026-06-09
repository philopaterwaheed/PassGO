package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const defaultRateLimitKeyPrefix = "passgo:rate_limit"

// RateLimitOptions configures IP-based request limiting backed by Redis.
type RateLimitOptions struct {
	RedisURL      string
	Requests      int
	WindowSeconds int
	KeyPrefix     string
}

// NewIPRateLimiter creates a Gin middleware that limits requests per client IP.
func NewIPRateLimiter(options RateLimitOptions) (gin.HandlerFunc, error) {
	if options.RedisURL == "" {
		return nil, nil
	}
	if options.Requests <= 0 {
		options.Requests = 100
	}
	if options.WindowSeconds <= 0 {
		options.WindowSeconds = 60
	}
	if options.KeyPrefix == "" {
		options.KeyPrefix = defaultRateLimitKeyPrefix
	}

	redisOptions, err := redis.ParseURL(options.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(redisOptions)
	window := time.Duration(options.WindowSeconds) * time.Second

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = "unknown"
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		key := fmt.Sprintf("%s:%s", options.KeyPrefix, ip)
		count, err := client.Incr(ctx, key).Result()
		if err != nil {
			log.Printf("Warning: rate limit check failed: %v", err)
			c.Next()
			return
		}

		if count == 1 {
			if err := client.Expire(ctx, key, window).Err(); err != nil {
				log.Printf("Warning: rate limit expiry failed: %v", err)
			}
		}

		ttl, err := client.TTL(ctx, key).Result()
		if err != nil || ttl < 0 {
			ttl = window
		}

		remaining := options.Requests - int(count)
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(options.Requests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))

		if count > int64(options.Requests) {
			c.Header("Retry-After", strconv.Itoa(int(ttl.Seconds())))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}, nil
}
