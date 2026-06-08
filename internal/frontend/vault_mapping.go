package frontend

import (
	"github.com/philopaterwaheed/passGO/internal/frontend/api"
	"github.com/philopaterwaheed/passGO/internal/frontend/state"
)

func apiVaultToState(v api.VaultResponse) state.Vault {
	return state.Vault{
		ID:          v.ID,
		Title:       v.Title,
		Username:    v.Username,
		Password:    v.Password,
		URL:         v.URL,
		Notes:       v.Notes,
		Decrypted:   v.Decrypted,
		HasPassword: v.HasPassword,
	}
}
