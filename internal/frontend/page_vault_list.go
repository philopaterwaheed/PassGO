package frontend

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

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
	invalidate func(),
) bool {
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
		invalidate()
		return true
	}

	return false
}
