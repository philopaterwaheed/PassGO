package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// User represents a user in the system
type User struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email         string        `bson:"email" json:"email" binding:"required,email"`
	SupabaseUID   string        `bson:"supabase_uid,omitempty" json:"supabase_uid,omitempty"`
	EmailVerified bool          `bson:"email_verified" json:"email_verified"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
	IsActive      bool          `bson:"is_active" json:"is_active"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// LoginRequest represents the login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the user data sent to clients (without sensitive info)
type UserResponse struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IsActive      bool      `json:"is_active"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID.Hex(),
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		IsActive:      u.IsActive,
	}
}

// SignupRequest represents the signup request
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// VerifyEmailRequest represents the email verification request
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Token string `json:"token" binding:"required"`
}

// ResendVerificationRequest represents the resend verification email request
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the password reset request
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// UpdatePasswordRequest represents the update password request
type UpdatePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// AuthResponse represents the authentication response with token
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
