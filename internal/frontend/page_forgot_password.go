package frontend

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/api"
	"github.com/philopaterwaheed/passGO/internal/frontend/pages"
	"github.com/philopaterwaheed/passGO/internal/frontend/state"
)

func handleForgotPasswordPage(
	gtx layout.Context,
	th *material.Theme,
	st *state.AppState,
	page *pages.ForgotPasswordPage,
	apiClient *api.Client,
	invalidate func(),
) {
	if page.BackBtn.Clicked(gtx) {
		st.Route = state.RouteLogin
		page.Reset()
	}

	if page.SendBtn.Clicked(gtx) && !page.IsLoading {
		email := page.EmailInput.Text()

		if email == "" {
			page.ErrorMsg = "Email is required"
		} else if !strings.Contains(email, "@") {
			page.ErrorMsg = "Invalid email address"
		} else {
			page.IsLoading = true
			page.ErrorMsg = ""
			page.SuccessMsg = ""

			go func() {
				msg, err := apiClient.ForgotPassword(email)
				if err != nil {
					page.ErrorMsg = err.Error()
					page.IsLoading = false
					invalidate()
					return
				}

				if msg == "" {
					msg = "If your email is registered, you will receive a reset link."
				}
				page.SuccessMsg = msg
				page.IsLoading = false
				invalidate()
			}()
		}
	}

	page.Layout(gtx, th)
}
