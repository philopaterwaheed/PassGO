package frontend

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

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
	invalidate func(),
) bool {
	shell.Layout(gtx, th, st, "Add vault", func(gtx layout.Context) layout.Dimensions {
		d, _ := page.Layout(gtx, th)
		return d
	})

	if page.BackBtn.Clicked(gtx) {
		st.Nav = state.NavVault
		st.Route = state.RouteVaultList
		invalidate()
		return true
	}

	if page.SaveBtn.Clicked(gtx) {
		v, ok := page.TryBuildVault()
		if !ok {
			invalidate()
			return true
		}
		st.AddVault(v)
		st.SelectedVaultID = st.Vaults[0].ID
		st.Nav = state.NavVault
		st.Route = state.RouteVaultDetail
		invalidate()
		return true
	}

	return false
}
