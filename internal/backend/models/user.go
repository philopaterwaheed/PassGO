package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// User represents a user in the system
type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string        `bson:"email" json:"email" binding:"required,email"`
	Username  string        `bson:"username" json:"username" binding:"required,min=3,max=50"`
	Password  string        `bson:"password" json:"-"`
	FullName  string        `bson:"full_name,omitempty" json:"full_name,omitempty"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
	IsActive  bool          `bson:"is_active" json:"is_active"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name,omitempty"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Username string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	FullName string `json:"full_name,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// LoginRequest represents the login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the user data sent to clients (without sensitive info)
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID.Hex(),
		Email:     u.Email,
		Username:  u.Username,
		FullName:  u.FullName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		IsActive:  u.IsActive,
	}
}
