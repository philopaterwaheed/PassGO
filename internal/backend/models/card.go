package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Card represents an encrypted credit card entry
type Card struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     string        `bson:"user_id" json:"user_id"`
	Title      string        `bson:"title" json:"title" binding:"required"`
	Cardholder string        `bson:"cardholder" json:"cardholder"`
	Number     string        `bson:"number" json:"number" binding:"required"`
	Expiry     string        `bson:"expiry" json:"expiry" binding:"required"`
	CVV        string        `bson:"cvv" json:"cvv" binding:"required"`
	Notes      string        `bson:"notes" json:"notes"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at" json:"updated_at"`
}
