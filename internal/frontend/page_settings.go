package frontend

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/pages"
	"github.com/philopaterwaheed/passGO/internal/frontend/state"
	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

func handleSettingsPage(
	gtx layout.Context,
	th *material.Theme,
	st *state.AppState,
	shell *ui.Shell,
	page *pages.SettingsPage,
	apiBaseURL string,
	invalidate func(),
) bool {
	var themeChanged bool
	shell.Layout(gtx, th, st, "Settings", func(gtx layout.Context) layout.Dimensions {
		var d layout.Dimensions
		d, themeChanged = page.Layout(gtx, th, st.Auth.Email, apiBaseURL, &st.DarkMode)
		return d
	})

	if themeChanged {
		invalidate()
		return true
	}

	return false
}
