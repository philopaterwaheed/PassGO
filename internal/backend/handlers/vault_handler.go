package handlers

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/philopaterwaheed/passGO/internal/backend/crypto"
	"github.com/philopaterwaheed/passGO/internal/backend/database"
	"github.com/philopaterwaheed/passGO/internal/backend/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type updateMasterPasswordRequest struct {
	CurrentMasterPassword string `json:"current_master_password"`
	NewMasterPassword     string `json:"new_master_password"`
}

// VaultHandler handles HTTP requests for vaults
type VaultHandler struct {
	repo     *database.VaultRepository
	userRepo *database.UserRepository
}

// NewVaultHandler creates a new VaultHandler
func NewVaultHandler() *VaultHandler {
	return &VaultHandler{
		repo:     database.NewVaultRepository(),
		userRepo: database.NewUserRepository(),
	}
}

func getMasterPassword(c *gin.Context) string {
	return c.GetHeader("X-Master-Password")
}

func (h *VaultHandler) getVaultKey(c *gin.Context, userID string) ([]byte, bool) {
	masterPassword := getMasterPassword(c)
	if masterPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Master-Password header is required"})
		return nil, false
	}

	user, err := h.userRepo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user vault key"})
		return nil, false
	}

	if user.MasterSalt == "" || user.VaultKey == "" {
		vaultKey, masterSalt, wrappedVaultKey, err := createWrappedVaultKey(masterPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vault key"})
			return nil, false
		}
		if err := h.userRepo.UpdateVaultKey(c.Request.Context(), userID, masterSalt, wrappedVaultKey); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save vault key"})
			return nil, false
		}
		return vaultKey, true
	}

	vaultKey, err := unwrapUserVaultKey(user, masterPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid master password"})
		return nil, false
	}

	return vaultKey, true
}

func createWrappedVaultKey(masterPassword string) ([]byte, string, string, error) {
	masterSalt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, "", "", err
	}
	vaultKey, err := crypto.GenerateVaultKey()
	if err != nil {
		return nil, "", "", err
	}
	masterKey := crypto.DeriveMasterKey(masterPassword, masterSalt)
	wrappedVaultKey, err := crypto.WrapVaultKey(vaultKey, masterKey)
	if err != nil {
		return nil, "", "", err
	}

	return vaultKey, base64.StdEncoding.EncodeToString(masterSalt), wrappedVaultKey, nil
}

func unwrapUserVaultKey(user *models.User, masterPassword string) ([]byte, error) {
	masterSalt, err := base64.StdEncoding.DecodeString(user.MasterSalt)
	if err != nil {
		return nil, err
	}

	masterKey := crypto.DeriveMasterKey(masterPassword, masterSalt)
	return crypto.UnwrapVaultKey(user.VaultKey, masterKey)
}

func wrapExistingVaultKey(vaultKey []byte, masterPassword string) (string, string, error) {
	masterSalt, err := crypto.GenerateSalt()
	if err != nil {
		return "", "", err
	}

	masterKey := crypto.DeriveMasterKey(masterPassword, masterSalt)
	wrappedVaultKey, err := crypto.WrapVaultKey(vaultKey, masterKey)
	if err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(masterSalt), wrappedVaultKey, nil
}

// UpdateMasterPassword re-wraps the vault key with a new master password.
func (h *VaultHandler) UpdateMasterPassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req updateMasterPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.NewMasterPassword == "" || len(req.NewMasterPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New master password must be at least 8 characters"})
		return
	}

	user, err := h.userRepo.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user vault key"})
		return
	}

	var vaultKey []byte
	if user.MasterSalt == "" || user.VaultKey == "" {
		vaultKey, err = crypto.GenerateVaultKey()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vault key"})
			return
		}
	} else {
		if req.CurrentMasterPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Current master password is required"})
			return
		}
		vaultKey, err = unwrapUserVaultKey(user, req.CurrentMasterPassword)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid current master password"})
			return
		}
	}

	masterSalt, wrappedVaultKey, err := wrapExistingVaultKey(vaultKey, req.NewMasterPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to protect vault key"})
		return
	}

	if err := h.userRepo.UpdateVaultKey(c.Request.Context(), userID.(string), masterSalt, wrappedVaultKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update master password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Master password updated successfully"})
}

// CreateVault handles POST requests to create a new vault
func (h *VaultHandler) CreateVault(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	vaultKey, ok := h.getVaultKey(c, userID.(string))
	if !ok {
		return
	}

	var req models.Vault
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Encrypt sensitive fields
	var err error
	if req.Password != "" {
		req.Password, err = crypto.Encrypt(req.Password, vaultKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
	}
	if req.Notes != "" {
		req.Notes, err = crypto.Encrypt(req.Notes, vaultKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt notes"})
			return
		}
	}

	req.UserID = userID.(string)

	if err := h.repo.CreateVault(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vault entry"})
		return
	}

	c.JSON(http.StatusCreated, req)
}

// GetVaults handles GET requests to retrieve all vaults for the current user
func (h *VaultHandler) GetVaults(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	vaultKey, ok := h.getVaultKey(c, userID.(string))
	if !ok {
		return
	}

	vaults, err := h.repo.GetVaultsByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vaults"})
		return
	}

	// Decrypt fields
	for i := range vaults {
		if vaults[i].Password != "" {
			dec, err := crypto.Decrypt(vaults[i].Password, vaultKey)
			if err == nil {
				vaults[i].Password = dec
			} else {
				log.Printf("Failed to decrypt password for vault %s\n", vaults[i].ID.Hex())
			}
		}
		if vaults[i].Notes != "" {
			dec, err := crypto.Decrypt(vaults[i].Notes, vaultKey)
			if err == nil {
				vaults[i].Notes = dec
			} else {
				log.Printf("Failed to decrypt notes for vault %s\n", vaults[i].ID.Hex())
			}
		}
	}

	c.JSON(http.StatusOK, vaults)
}

// GetVault handles GET requests to retrieve a single vault
func (h *VaultHandler) GetVault(c *gin.Context) {
	idParam := c.Param("id")
	id, err := bson.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vault ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	vault, err := h.repo.GetVaultByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vault not found"})
		return
	}

	if vault.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	vaultKey, ok := h.getVaultKey(c, userID.(string))
	if !ok {
		return
	}

	if vault.Password != "" {
		dec, err := crypto.Decrypt(vault.Password, vaultKey)
		if err == nil {
			vault.Password = dec
		}
	}
	if vault.Notes != "" {
		dec, err := crypto.Decrypt(vault.Notes, vaultKey)
		if err == nil {
			vault.Notes = dec
		}
	}

	c.JSON(http.StatusOK, vault)
}

// UpdateVault handles PUT requests to update a vault
func (h *VaultHandler) UpdateVault(c *gin.Context) {
	idParam := c.Param("id")
	id, err := bson.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vault ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	existingVault, err := h.repo.GetVaultByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vault not found"})
		return
	}

	if existingVault.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	vaultKey, ok := h.getVaultKey(c, userID.(string))
	if !ok {
		return
	}

	var req models.Vault
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Encrypt sensitive fields
	if req.Password != "" {
		req.Password, err = crypto.Encrypt(req.Password, vaultKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
	}
	if req.Notes != "" {
		req.Notes, err = crypto.Encrypt(req.Notes, vaultKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt notes"})
			return
		}
	}

	req.UserID = userID.(string)

	if err := h.repo.UpdateVault(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vault"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vault updated successfully"})
}

// DeleteVault handles DELETE requests to a vault
func (h *VaultHandler) DeleteVault(c *gin.Context) {
	idParam := c.Param("id")
	id, err := bson.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vault ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	existingVault, err := h.repo.GetVaultByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vault not found"})
		return
	}

	if existingVault.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.repo.DeleteVault(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vault"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vault deleted successfully"})
}