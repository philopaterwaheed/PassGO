package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Vault represents an encrypted password vault or credit card entry
type Vault struct {
	ID     bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID string        `bson:"user_id" json:"user_id"`
	Type   string        `bson:"type" json:"type"` // "login" or "card"
	Title  string        `bson:"title" json:"title" binding:"required"`
	// Fields for "login" type
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
	URL      string `bson:"url" json:"url"`
	// Fields for "card" type
	Cardholder string `bson:"cardholder" json:"cardholder"`
	Number     string `bson:"number" json:"number"`
	Expiry     string `bson:"expiry" json:"expiry"`
	CVV        string `bson:"cvv" json:"cvv"`
	// Common fields
	Notes     string    `bson:"notes" json:"notes"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	Decrypted bool      `bson:"-" json:"decrypted"`
}
