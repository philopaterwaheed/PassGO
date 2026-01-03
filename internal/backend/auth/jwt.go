package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/philopaterwaheed/passGO/internal/backend/config"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
	ErrTokenClaims  = errors.New("invalid token claims")
)

// Claims represents the JWT claims
type Claims struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	SupabaseUID string `json:"supabase_uid"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for a user
func GenerateToken(userID, email, supabaseUID string) (string, error) {
	if config.JWTSecret == "" {
		return "", errors.New("JWT secret not configured")
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:      userID,
		Email:       email,
		SupabaseUID: supabaseUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "passgo-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken validates a JWT token and returns the claims
func VerifyToken(tokenString string) (*Claims, error) {
	if config.JWTSecret == "" {
		return nil, errors.New("JWT secret not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenClaims
	}

	return claims, nil
}

// RefreshToken creates a new token with extended expiration
func RefreshToken(oldToken string) (string, error) {
	claims, err := VerifyToken(oldToken)
	if err != nil {
		// Allow expired tokens to be refreshed within 7 days
		if !errors.Is(err, ErrTokenExpired) {
			return "", err
		}

		token, parseErr := jwt.ParseWithClaims(oldToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		}, jwt.WithoutClaimsValidation())

		if parseErr != nil {
			return "", ErrInvalidToken
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return "", ErrTokenClaims
		}

		if time.Since(claims.IssuedAt.Time) > 7*24*time.Hour {
			return "", ErrTokenExpired
		}
	}

	return GenerateToken(claims.UserID, claims.Email, claims.SupabaseUID)
}
