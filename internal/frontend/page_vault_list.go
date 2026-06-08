package frontend

import (
	"log"
	"strings"

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
	if st.VaultsLoadDone != nil {
		select {
		case result := <-st.VaultsLoadDone:
			st.VaultsLoading = false
			st.VaultsLoadDone = nil
			if result.Err != nil {
				log.Printf("Failed to fetch vaults: %v", result.Err)
				st.VaultsLoadError = result.Err.Error()
				if strings.Contains(strings.ToLower(result.Err.Error()), "master password") {
					st.Auth.MasterPassword = ""
					page.UnlockError = result.Err.Error()
					page.MasterPasswordInput.SetText("")
					st.VaultsLoadError = ""
				}
				st.VaultsLoaded = true
			} else {
				st.Vaults = result.Vaults
				st.VaultsLoadError = ""
				st.VaultsLoaded = true
				page.UnlockError = ""
				page.MasterPasswordInput.SetText("")
			}
		default:
		}
	}

	// Automatically fetch vaults once the vault has been unlocked.
	if !st.VaultsLoaded && !st.VaultsLoading && st.Auth.MasterPassword != "" {
		st.VaultsLoading = true
		st.VaultsLoadError = ""
		masterPassword := st.Auth.MasterPassword
		done := make(chan state.VaultLoadResult, 1)
		st.VaultsLoadDone = done

		go func() {
			resp, err := apiClient.GetVaults(masterPassword)
			if err != nil {
				done <- state.VaultLoadResult{Err: err}
				invalidate()
				return
			}

			// Map to state models
			var newVaults []state.Vault
			for _, v := range resp {
				newVaults = append(newVaults, apiVaultToState(v))
			}
			done <- state.VaultLoadResult{Vaults: newVaults}
			invalidate()
		}()
	}

	var action pages.VaultListAction
	shell.Layout(gtx, th, st, "Vaults", func(gtx layout.Context) layout.Dimensions {
		var d layout.Dimensions
		d, action = page.Layout(gtx, th, st.Vaults, st.VaultsLoading, st.VaultsLoadError, st.Auth.MasterPassword == "")
		return d
	})

	if action.Unlock {
		if action.UnlockPassword == "" {
			page.UnlockError = "Master password is required"
			invalidate()
			return true
		}
		st.Auth.MasterPassword = action.UnlockPassword
		page.UnlockError = ""
		st.VaultsLoaded = false
		st.VaultsLoading = false
		st.VaultsLoadError = ""
		st.VaultsLoadDone = nil
		invalidate()
		return true
	}

	if action.Retry {
		st.VaultsLoaded = false
		st.VaultsLoading = false
		st.VaultsLoadError = ""
		st.VaultsLoadDone = nil
		invalidate()
		return true
	}

	if action.Add {
		st.Nav = state.NavVault
		st.Route = state.RouteVaultAdd
		vaultAddPage.Reset(nil)
		invalidate()
		return true
	}

	if action.DeleteID != "" {
		go func(id string) {
			err := apiClient.DeleteVault(id)
			if err != nil {
				log.Printf("Failed to delete vault: %v", err)
				invalidate()
				return
			}
			st.DeleteVault(id)
			invalidate()
		}(action.DeleteID)
		return true
	}

	if action.EditID != "" {
		st.SelectedVaultID = action.EditID
		// finding the vault to pass to reset
		var vToEdit *state.Vault
		for _, v := range st.Vaults {
			if v.ID == action.EditID {
				vToEdit = &v
				break
			}
		}
		if vToEdit != nil && !vToEdit.Decrypted {
			go func(id string, masterPassword string) {
				resp, err := apiClient.GetVault(id, masterPassword)
				if err != nil {
					log.Printf("Failed to fetch vault before edit: %v", err)
					invalidate()
					return
				}

				v := apiVaultToState(*resp)
				st.UpdateVault(v)
				st.Nav = state.NavVault
				st.Route = state.RouteVaultAdd
				vaultAddPage.Reset(&v)
				invalidate()
			}(action.EditID, st.Auth.MasterPassword)
			return true
		}
		st.Nav = state.NavVault
		st.Route = state.RouteVaultAdd
		vaultAddPage.Reset(vToEdit)
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

	return st.VaultsLoading
}
