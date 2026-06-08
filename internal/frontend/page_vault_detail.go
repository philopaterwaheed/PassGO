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

func handleVaultDetailPage(
	gtx layout.Context,
	th *material.Theme,
	st *state.AppState,
	shell *ui.Shell,
	page *pages.VaultDetailPage,
	apiClient *api.Client,
	invalidate func(),
) bool {
	if st.VaultDetailLoadDone != nil {
		select {
		case result := <-st.VaultDetailLoadDone:
			st.VaultDetailLoading = false
			st.VaultDetailLoadDone = nil
			if result.Err != nil {
				log.Printf("Failed to fetch vault detail: %v", result.Err)
				st.VaultDetailLoadError = result.Err.Error()
				if strings.Contains(strings.ToLower(result.Err.Error()), "master password") {
					st.Auth.MasterPassword = ""
				}
			} else {
				st.UpdateVault(result.Vault)
				st.VaultDetailLoadError = ""
			}
		default:
		}
	}

	v, ok := st.VaultByID(st.SelectedVaultID)
	if !ok {
		st.Route = state.RouteVaultList
		return true
	}

	if !v.Decrypted && !st.VaultDetailLoading && st.Auth.MasterPassword != "" {
		st.VaultDetailLoading = true
		st.VaultDetailLoadError = ""
		id := v.ID
		masterPassword := st.Auth.MasterPassword
		done := make(chan state.VaultDetailLoadResult, 1)
		st.VaultDetailLoadDone = done

		go func() {
			resp, err := apiClient.GetVault(id, masterPassword)
			if err != nil {
				done <- state.VaultDetailLoadResult{Err: err}
				invalidate()
				return
			}
			done <- state.VaultDetailLoadResult{Vault: apiVaultToState(*resp)}
			invalidate()
		}()
	}

	var action pages.VaultDetailAction
	shell.Layout(gtx, th, st, "Vault", func(gtx layout.Context) layout.Dimensions {
		var d layout.Dimensions
		d, action = page.Layout(gtx, th, v, st.VaultDetailLoading, st.VaultDetailLoadError)
		return d
	})

	if action.Back {
		st.Nav = state.NavVault
		st.Route = state.RouteVaultList
		st.VaultDetailLoading = false
		st.VaultDetailLoadError = ""
		st.VaultDetailLoadDone = nil
		invalidate()
		return true
	}

	return st.VaultDetailLoading
}
