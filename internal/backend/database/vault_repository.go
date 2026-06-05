package database

import (
	"context"
	"fmt"
	"time"

	"github.com/philopaterwaheed/passGO/internal/backend/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const vaultCollectionName = "vaults"

// VaultRepository handles database operations for vaults
type VaultRepository struct {
	collection *mongo.Collection
}

// NewVaultRepository creates a new vault repository
func NewVaultRepository() *VaultRepository {
	return &VaultRepository{
		collection: GetCollection(vaultCollectionName),
	}
}

// CreateVault inserts a new vault into the database
func (r *VaultRepository) CreateVault(ctx context.Context, vault *models.Vault) error {
	vault.ID = bson.NewObjectID()
	vault.CreatedAt = time.Now()
	vault.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, vault)
	return err
}

// GetVaultsByUserID retrieves all vaults for a specific user
func (r *VaultRepository) GetVaultsByUserID(ctx context.Context, userID string) ([]models.Vault, error) {
	var vaults []models.Vault
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &vaults); err != nil {
		return nil, err
	}

	if vaults == nil {
		vaults = []models.Vault{}
	}

	return vaults, nil
}

// GetVaultByID retrieves a specific vault by its ID
func (r *VaultRepository) GetVaultByID(ctx context.Context, id bson.ObjectID) (*models.Vault, error) {
	var vault models.Vault
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&vault)
	if err != nil {
		return nil, err
	}
	return &vault, nil
}

// UpdateVault updates an existing vault in the database
func (r *VaultRepository) UpdateVault(ctx context.Context, id bson.ObjectID, vault *models.Vault) error {
	vault.UpdatedAt = time.Now()
	
	update := bson.M{
		"$set": bson.M{
			"title":      vault.Title,
			"username":   vault.Username,
			"password":   vault.Password,
			"url":        vault.URL,
			"notes":      vault.Notes,
			"updated_at": vault.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("vault not found")
	}
	return nil
}

// DeleteVault deletes a vault from the database
func (r *VaultRepository) DeleteVault(ctx context.Context, id bson.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("vault not found")
	}
	return nil
}