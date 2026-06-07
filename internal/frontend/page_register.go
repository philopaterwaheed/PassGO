package frontend

import (
	"log"
	"strings"

	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/api"
	"github.com/philopaterwaheed/passGO/internal/frontend/pages"
	"github.com/philopaterwaheed/passGO/internal/frontend/state"
)

func handleRegisterPage(
	gtx layout.Context,
	th *material.Theme,
	st *state.AppState,
	page *pages.RegisterPage,
	apiClient *api.Client,
	invalidate func(),
) {
	if page.BackBtn.Clicked(gtx) {
		st.Route = state.RouteWelcome
		page.Reset()
	}

	if page.RegisterBtn.Clicked(gtx) && !page.IsLoading {
		email := page.EmailInput.Text()
		password := page.PasswordInput.Text()
		confirmPassword := page.ConfirmPasswordInput.Text()
		masterPassword := page.MasterPasswordInput.Text()
		confirmMasterPassword := page.ConfirmMasterPasswordInput.Text()

		if email == "" || password == "" || confirmPassword == "" || masterPassword == "" || confirmMasterPassword == "" {
			page.ErrorMsg = "All fields are required"
		} else if !strings.Contains(email, "@") {
			page.ErrorMsg = "Invalid email address"
		} else if len(password) < 8 {
			page.ErrorMsg = "Password must be at least 8 characters"
		} else if password != confirmPassword {
			page.ErrorMsg = "Passwords do not match"
		} else if len(masterPassword) < 8 {
			page.ErrorMsg = "Master password must be at least 8 characters"
		} else if masterPassword != confirmMasterPassword {
			page.ErrorMsg = "Master passwords do not match"
		} else {
			page.IsLoading = true
			page.ErrorMsg = ""
			page.SuccessMsg = ""

			go func() {
				resp, err := apiClient.Signup(email, password, masterPassword)
				if err != nil {
					page.ErrorMsg = err.Error()
					page.IsLoading = false
					invalidate()
					return
				}

				page.SuccessMsg = resp.Message
				if page.SuccessMsg == "" {
					page.SuccessMsg = "Registration successful! Please check your email."
				}
				page.IsLoading = false
				page.PasswordInput.SetText("")
				page.ConfirmPasswordInput.SetText("")
				page.MasterPasswordInput.SetText("")
				page.ConfirmMasterPasswordInput.SetText("")
				log.Printf("Registered successfully: %+v", resp.User)
				invalidate()
			}()
		}
	}

	page.Layout(gtx, th)
}
