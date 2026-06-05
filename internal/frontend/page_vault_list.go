package frontend

import (
	"log"

	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/api"
	"github.com/philopaterwaheed/passGO/internal/frontend/pages"
	"github.com/philopaterwaheed/passGO/internal/frontend/state"
	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

func handleVaultListPage(
	gtx layout.Context,
	th *material.Theme,
	st *state.AppState,
	shell *ui.Shell,
	page *pages.VaultListPage,
	vaultAddPage *pages.VaultAddPage,
	vaultDetailPage *pages.VaultDetailPage,
	apiClient *api.Client,
	invalidate func(),
) bool {
	// Automatically fetch vaults once
	if !st.VaultsLoaded && st.Auth.MasterPassword != "" {
		st.VaultsLoaded = true
		go func() {
			resp, err := apiClient.GetVaults(st.Auth.MasterPassword)
			if err != nil {
				log.Printf("Failed to fetch vaults: %v", err)
				// Might want to clear loaded flag to retry or show error
				return
			}
			
			// Map to state models
			var newVaults []state.Vault
			for _, v := range resp {
				newVaults = append(newVaults, state.Vault{
					ID:       v.ID,
					Title:    v.Title,
					Username: v.Username,
					Password: v.Password,
					URL:      v.URL,
					Notes:    v.Notes,
				})
			}
			st.Vaults = newVaults
			invalidate()
		}()
	}

	var action pages.VaultListAction
	shell.Layout(gtx, th, st, "Vaults", func(gtx layout.Context) layout.Dimensions {
		var d layout.Dimensions
		d, action = page.Layout(gtx, th, st.Vaults)
		return d
	})

	if action.Add {
		st.Nav = state.NavVault
		st.Route = state.RouteVaultAdd
		vaultAddPage.Reset()
		invalidate()
		return true
	}

	if action.OpenID != "" {
		st.SelectedVaultID = action.OpenID
		st.Nav = state.NavVault
		st.Route = state.RouteVaultDetail
		vaultDetailPage.Reset()
		invalidate()
		return true
	}

	return false
}
