package frontend

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

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
	invalidate func(),
) bool {
	v, ok := st.VaultByID(st.SelectedVaultID)
	if !ok {
		st.Route = state.RouteVaultList
		return true
	}

	var action pages.VaultDetailAction
	shell.Layout(gtx, th, st, "Vault", func(gtx layout.Context) layout.Dimensions {
		var d layout.Dimensions
		d, action = page.Layout(gtx, th, v)
		return d
	})

	if action.Back {
		st.Nav = state.NavVault
		st.Route = state.RouteVaultList
		invalidate()
		return true
	}

	return false
}
