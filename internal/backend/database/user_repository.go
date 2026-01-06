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
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("email already exists")
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

// GetAllUsers retrieves all users with pagination
func (r *UserRepository) GetAllUsers(ctx context.Context, page, limit int64) ([]*models.User, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates a user's information
func (r *UserRepository) UpdateUser(ctx context.Context, id string, update *models.UpdateUserRequest) (*models.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	updateDoc := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	setFields := updateDoc["$set"].(bson.M)

	if update.Email != "" {
		setFields["email"] = update.Email
	}
	if update.IsActive != nil {
		setFields["is_active"] = *update.IsActive
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser models.User
	err = r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		updateDoc,
		opts,
	).Decode(&updatedUser)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		if mongo.IsDuplicateKeyError(err) {
			if containsField(err.Error(), "email") {
				return nil, ErrDuplicateEmail
			}
		}
		return nil, err
	}

	return &updatedUser, nil
}

// DeleteUser deletes a user from the database
func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

// CountUsers returns the total number of users
func (r *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
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
