package database

import (
	"context"
	"errors"
	"time"

	"github.com/philopaterwaheed/passGO/internal/backend/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const usersCollection = "users"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateEmail = errors.New("email already exists")
)

// UserRepository handles user database operations
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		collection: GetCollection(usersCollection),
	}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	user.ID = bson.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// Check which field caused the duplicate key error
			if containsField(err.Error(), "email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetUserBySupabaseUID retrieves a user by their Supabase UID
func (r *UserRepository) GetUserBySupabaseUID(ctx context.Context, supabaseUID string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"supabase_uid": supabaseUID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// UpdateEmailVerified updates the email verification status
func (r *UserRepository) UpdateEmailVerified(ctx context.Context, id string, verified bool) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"email_verified": verified,
			"updated_at":     time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateVaultKey stores the salt and encrypted vault key for a user.
func (r *UserRepository) UpdateVaultKey(ctx context.Context, id, masterSalt, vaultKey string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"master_salt": masterSalt,
			"vault_key":   vaultKey,
			"updated_at":  time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

// CreateIndexes creates necessary indexes for the users collection
func (r *UserRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "supabase_uid", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Helper function to check if error message contains a field name
func containsField(errMsg, field string) bool {
	return len(errMsg) > 0 && len(field) > 0 &&
		(errMsg[0:1] == field[0:1] || errMsg[len(errMsg)-len(field):] == field)
}
