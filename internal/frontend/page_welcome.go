package frontend

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/pages"
	"github.com/philopaterwaheed/passGO/internal/frontend/state"
)

func handleWelcomePage(gtx layout.Context, th *material.Theme, st *state.AppState, page *pages.WelcomePage) {
	if page.LoginBtn.Clicked(gtx) {
		st.Route = state.RouteLogin
	}
	if page.RegisterBtn.Clicked(gtx) {
		st.Route = state.RouteRegister
	}

	page.Layout(gtx, th)
}
