package frontend

import (
	"log"

	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/api"
	"github.com/philopaterwaheed/passGO/internal/frontend/pages"
	"github.com/philopaterwaheed/passGO/internal/frontend/state"
	"github.com/philopaterwaheed/passGO/internal/frontend/storage"
)

func handleLoginPage(
	gtx layout.Context,
	th *material.Theme,
	st *state.AppState,
	page *pages.LoginPage,
	apiClient *api.Client,
	sessionStore storage.SessionStore,
	invalidate func(),
) {
	if page.BackBtn.Clicked(gtx) {
		st.Route = state.RouteWelcome
		page.Reset()
	}

	if page.ForgotBtn.Clicked(gtx) {
		st.Route = state.RouteForgotPassword
		page.ErrorMsg = ""
		page.SuccessMsg = ""
	}

	if page.LoginBtn.Clicked(gtx) && !page.IsLoading {
		email := page.EmailInput.Text()
		password := page.PasswordInput.Text()

		if email == "" || password == "" {
			page.ErrorMsg = "Email and password are required"
		} else {
			page.IsLoading = true
			page.ErrorMsg = ""
			page.SuccessMsg = ""

			go func() {
				resp, err := apiClient.Login(email, password)
				if err != nil {
					page.ErrorMsg = err.Error()
					page.IsLoading = false
					invalidate()
					return
				}

				st.Auth.Token = resp.Token
				st.Auth.Email = resp.User.Email
				st.Auth.MasterPassword = password
				st.Vaults = nil
				st.VaultsLoaded = false
				st.VaultsLoading = false
				st.VaultsLoadError = ""
				st.VaultsLoadDone = nil
				_ = sessionStore.Save(storage.Session{Token: resp.Token, Email: resp.User.Email, MasterPassword: password})
				page.SuccessMsg = "Welcome, " + resp.User.Email
				page.IsLoading = false
				st.Nav = state.NavVault
				st.Route = state.RouteVaultList
				log.Printf("Logged in successfully: %+v", resp.User)
				invalidate()
			}()
		}
	}

	page.Layout(gtx, th)
}
