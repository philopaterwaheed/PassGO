package frontend

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/api"
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
	apiClient *api.Client,
	apiBaseURL string,
	invalidate func(),
) bool {
	var themeChanged bool
	var action pages.SettingsAction
	shell.Layout(gtx, th, st, "Settings", func(gtx layout.Context) layout.Dimensions {
		var d layout.Dimensions
		d, themeChanged, action = page.Layout(gtx, th, st.Auth.Email, apiBaseURL, &st.DarkMode)
		return d
	})

	if action.NewMasterPassword != "" || action.CurrentMasterPassword != "" || action.ConfirmMasterPassword != "" {
		action.UpdateMasterPassword = true
	}

	if action.UpdateMasterPassword {
		if page.IsSaving {
			return false
		}
		if action.CurrentMasterPassword == "" {
			page.ErrorMsg = "Current master password is required"
			page.SuccessMsg = ""
			invalidate()
			return true
		}
		if len(action.NewMasterPassword) < 8 {
			page.ErrorMsg = "New master password must be at least 8 characters"
			page.SuccessMsg = ""
			invalidate()
			return true
		}
		if action.NewMasterPassword != action.ConfirmMasterPassword {
			page.ErrorMsg = "New master passwords do not match"
			page.SuccessMsg = ""
			invalidate()
			return true
		}

		page.IsSaving = true
		page.ErrorMsg = ""
		page.SuccessMsg = ""
		invalidate()

		currentMasterPassword := action.CurrentMasterPassword
		newMasterPassword := action.NewMasterPassword

		go func() {
			err := apiClient.UpdateMasterPassword(currentMasterPassword, newMasterPassword)
			if err != nil {
				page.ErrorMsg = err.Error()
				page.IsSaving = false
				invalidate()
				return
			}

			st.Auth.MasterPassword = newMasterPassword
			page.ClearMasterPasswordFields()
			page.SuccessMsg = "Master password updated"
			page.IsSaving = false
			invalidate()
		}()

		return true
	}

	if action.UpdateAccountPassword {
		if page.IsSaving {
			return false
		}
		if action.CurrentAccountPassword == "" {
			page.ErrorMsg = "Current account password is required"
			page.SuccessMsg = ""
			invalidate()
			return true
		}
		if len(action.NewAccountPassword) < 8 {
			page.ErrorMsg = "New account password must be at least 8 characters"
			page.SuccessMsg = ""
			invalidate()
			return true
		}
		if action.NewAccountPassword != action.ConfirmAccountPassword {
			page.ErrorMsg = "New account passwords do not match"
			page.SuccessMsg = ""
			invalidate()
			return true
		}

		page.IsSaving = true
		page.ErrorMsg = ""
		page.SuccessMsg = ""
		invalidate()

		currentAccountPassword := action.CurrentAccountPassword
		newAccountPassword := action.NewAccountPassword

		go func() {
			err := apiClient.UpdateAccountPassword(currentAccountPassword, newAccountPassword)
			if err != nil {
				page.ErrorMsg = err.Error()
				page.IsSaving = false
				invalidate()
				return
			}

			page.ClearAccountPasswordFields()
			page.SuccessMsg = "Account password updated"
			page.IsSaving = false
			invalidate()
		}()

		return true
	}

	if themeChanged {
		invalidate()
		return true
	}

	return false
}
