package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/philopaterwaheed/passGO/internal/backend/auth"
	"github.com/philopaterwaheed/passGO/internal/backend/database"
	"github.com/philopaterwaheed/passGO/internal/backend/models"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	repo     *database.UserRepository
	supabase *auth.SupabaseClient
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() (*AuthHandler, error) {
	supabaseClient, err := auth.NewSupabaseClient()
	if err != nil {
		return nil, err
	}

	return &AuthHandler{
		repo:     database.NewUserRepository(),
		supabase: supabaseClient,
	}, nil
}

// Signup handles POST /api/auth/signup
func (h *AuthHandler) Signup(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists in local database
	_, err := h.repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Register with Supabase
	supabaseResp, err := h.supabase.SignUp(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists in authentication system"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register: " + err.Error()})
		return
	}

	// Create user in local database
	user := &models.User{
		Email:         req.Email,
		SupabaseUID:   supabaseResp.User.ID,
		EmailVerified: false,
	}

	if err := h.repo.CreateUser(c.Request.Context(), user); err != nil {
		if errors.Is(err, database.ErrDuplicateEmail) {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		if errors.Is(err, database.ErrDuplicateUsername) {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please check your email to verify your account.",
		"user":    user.ToResponse(),
	})
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate with Supabase
	supabaseResp, err := h.supabase.SignIn(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		if errors.Is(err, auth.ErrEmailNotVerified) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your email before logging in"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	user, err := h.repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			user = &models.User{
				Email:         supabaseResp.User.Email,
				SupabaseUID:   supabaseResp.User.ID,
				EmailVerified: true,
				IsActive:      true,
			}
			if err := h.repo.CreateUser(c.Request.Context(), user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync user"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}
	}

	// Update email verification status if needed
	if !user.EmailVerified && supabaseResp.User.EmailConfirmedAt != "" {
		if err := h.repo.UpdateEmailVerified(c.Request.Context(), user.ID.Hex(), true); err == nil {
			user.EmailVerified = true
		}
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is disabled"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID.Hex(), user.Email, user.SupabaseUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

// VerifyEmail handles POST /api/auth/verify-email
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req models.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	supabaseResp, err := h.supabase.VerifyOTP(req.Email, req.Token, "signup")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification token"})
		return
	}

	user, err := h.repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := h.repo.UpdateEmailVerified(c.Request.Context(), user.ID.Hex(), true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verification status"})
		return
	}

	token, err := auth.GenerateToken(user.ID.Hex(), user.Email, supabaseResp.User.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	user.EmailVerified = true

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
		"token":   token,
		"user":    user.ToResponse(),
	})
}

// ResendVerification handles POST /api/auth/resend-verification
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req models.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	user, err := h.repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			// Don't reveal if email exists for security
			c.JSON(http.StatusOK, gin.H{"message": "If your email is registered, you will receive a verification email"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	// Check if already verified
	if user.EmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already verified"})
		return
	}

	// Resend verification email via Supabase
	// Todo add rate limiting to prevent abuse
	if err := h.supabase.ResendVerificationEmail(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

// ForgotPassword handles POST /api/auth/forgot-password
// Sends a password reset email
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send password reset email via Supabase
	// Don't check if user exists for security reasons
	if err := h.supabase.ResetPassword(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send password reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If your email is registered, you will receive a password reset email"})
}

// RefreshToken handles POST /api/auth/refresh
// Refreshes an existing JWT token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get current token from header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Refresh the token
	newToken, err := auth.RefreshToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

// GetCurrentUser handles GET /api/auth/me
// Returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.repo.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

