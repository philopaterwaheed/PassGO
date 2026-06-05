package frontend

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/api"
	"github.com/philopaterwaheed/passGO/internal/frontend/pages"
	"github.com/philopaterwaheed/passGO/internal/frontend/state"
	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

func handleVaultAddPage(
	gtx layout.Context,
	th *material.Theme,
	st *state.AppState,
	shell *ui.Shell,
	page *pages.VaultAddPage,
	apiClient *api.Client,
	invalidate func(),
) bool {
	var action pages.VaultAddAction
	shell.Layout(gtx, th, st, "Add vault", func(gtx layout.Context) layout.Dimensions {
		var d layout.Dimensions
		d, action = page.Layout(gtx, th)
		return d
	})

	if action.Back {
		st.Nav = state.NavVault
		st.Route = state.RouteVaultList
		invalidate()
		return true
	}

	if action.Save && !page.IsLoading {
		page.IsLoading = true
		page.ErrorMsg = ""
		invalidate()

		v := action.Vault
		go func() {
			resp, err := apiClient.CreateVault(st.Auth.MasterPassword, &api.VaultRequest{
				Title:    v.Title,
				Username: v.Username,
				Password: v.Password,
				URL:      v.URL,
				Notes:    v.Notes,
			})

			if err != nil {
				page.IsLoading = false
				page.ErrorMsg = err.Error()
				invalidate()
				return
			}

			// Map response back to state vault
			v.ID = resp.ID
			st.AddVault(v)

			page.IsLoading = false
			st.SelectedVaultID = v.ID
			st.Nav = state.NavVault
			st.Route = state.RouteVaultDetail
			invalidate()
		}()

		return true
	}

	return false
}
