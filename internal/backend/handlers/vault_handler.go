package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/philopaterwaheed/passGO/internal/backend/crypto"
	"github.com/philopaterwaheed/passGO/internal/backend/database"
	"github.com/philopaterwaheed/passGO/internal/backend/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// VaultHandler handles HTTP requests for vaults
type VaultHandler struct {
	repo *database.VaultRepository
}

// NewVaultHandler creates a new VaultHandler
func NewVaultHandler() *VaultHandler {
	return &VaultHandler{
		repo: database.NewVaultRepository(),
	}
}

func getMasterPassword(c *gin.Context) string {
	return c.GetHeader("X-Master-Password")
}

// CreateVault handles POST requests to create a new vault
func (h *VaultHandler) CreateVault(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	masterPassword := getMasterPassword(c)
	if masterPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Master-Password header is required"})
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
		req.Password, err = crypto.Encrypt(req.Password, masterPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
	}
	if req.Notes != "" {
		req.Notes, err = crypto.Encrypt(req.Notes, masterPassword)
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

	masterPassword := getMasterPassword(c)
	if masterPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Master-Password header is required"})
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
			dec, err := crypto.Decrypt(vaults[i].Password, masterPassword)
			if err == nil {
				vaults[i].Password = dec
			} else {
				log.Printf("Failed to decrypt password for vault %s\n", vaults[i].ID.Hex())
			}
		}
		if vaults[i].Notes != "" {
			dec, err := crypto.Decrypt(vaults[i].Notes, masterPassword)
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

	masterPassword := getMasterPassword(c)
	if masterPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Master-Password header is required"})
		return
	}

	if vault.Password != "" {
		dec, err := crypto.Decrypt(vault.Password, masterPassword)
		if err == nil {
			vault.Password = dec
		}
	}
	if vault.Notes != "" {
		dec, err := crypto.Decrypt(vault.Notes, masterPassword)
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

	masterPassword := getMasterPassword(c)
	if masterPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Master-Password header is required"})
		return
	}

	var req models.Vault
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Encrypt sensitive fields
	if req.Password != "" {
		req.Password, err = crypto.Encrypt(req.Password, masterPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
	}
	if req.Notes != "" {
		req.Notes, err = crypto.Encrypt(req.Notes, masterPassword)
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