package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

type SettingsPage struct {
	DarkMode widget.Bool

	CurrentMasterPasswordInput widget.Editor
	NewMasterPasswordInput     widget.Editor
	ConfirmMasterPasswordInput widget.Editor
	SaveMasterPasswordBtn      widget.Clickable
	GitHubBtn                  widget.Clickable
	LinkedInBtn                widget.Clickable
	PortfolioBtn               widget.Clickable
	CVBtn                      widget.Clickable
	ErrorMsg                   string
	SuccessMsg                 string
	IsSaving                   bool

	initialized bool
}

type SettingsAction struct {
	CurrentMasterPassword string
	NewMasterPassword     string
	ConfirmMasterPassword string
}

func NewSettingsPage() *SettingsPage {
	return &SettingsPage{
		CurrentMasterPasswordInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
			Mask:       '*',
		},
		NewMasterPasswordInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
			Mask:       '*',
		},
		ConfirmMasterPasswordInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
			Mask:       '*',
		},
	}
}

func (p *SettingsPage) ClearMasterPasswordFields() {
	p.CurrentMasterPasswordInput.SetText("")
	p.NewMasterPasswordInput.SetText("")
	p.ConfirmMasterPasswordInput.SetText("")
}

func (p *SettingsPage) Reset() {
	p.ClearMasterPasswordFields()
	p.ErrorMsg = ""
	p.SuccessMsg = ""
	p.IsSaving = false
}

func (p *SettingsPage) Layout(gtx layout.Context, th *material.Theme, email string, apiBaseURL string, darkMode *bool) (layout.Dimensions, bool, SettingsAction) {
	if darkMode != nil && !p.initialized {
		p.DarkMode.Value = *darkMode
		p.initialized = true
	}

	changed := false
	var action SettingsAction

	if p.SaveMasterPasswordBtn.Clicked(gtx) {
		action.CurrentMasterPassword = p.CurrentMasterPasswordInput.Text()
		action.NewMasterPassword = p.NewMasterPasswordInput.Text()
		action.ConfirmMasterPassword = p.ConfirmMasterPasswordInput.Text()
	}
	for p.GitHubBtn.Clicked(gtx) {
		ui.OpenURLLogged(ui.GitHubURL)
	}
	for p.LinkedInBtn.Clicked(gtx) {
		ui.OpenURLLogged(ui.LinkedInURL)
	}
	for p.PortfolioBtn.Clicked(gtx) {
		ui.OpenURLLogged(ui.PortfolioURL)
	}
	for p.CVBtn.Clicked(gtx) {
		ui.OpenURLLogged(ui.CVURL)
	}

	dims := layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(640)))
		return ui.SurfacePanel(gtx, ui.SurfaceColor, ui.RadiusLarge, layout.UniformInset(unit.Dp(24)), func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.H6(th, "Account")
					lbl.Color = th.Palette.ContrastBg
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(ui.MutedLabel(th, material.Body2(th, "Signed in as: "+email)).Layout),
				layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
				layout.Rigid(material.Body1(th, "Master password").Layout),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					e := material.Editor(th, &p.CurrentMasterPasswordInput, "Current master password")
					e.TextSize = unit.Sp(16)
					return e.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					e := material.Editor(th, &p.NewMasterPasswordInput, "New master password")
					e.TextSize = unit.Sp(16)
					return e.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					e := material.Editor(th, &p.ConfirmMasterPasswordInput, "Confirm new master password")
					e.TextSize = unit.Sp(16)
					return e.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if p.ErrorMsg != "" {
						lbl := material.Body2(th, p.ErrorMsg)
						lbl.Color = ui.DangerColor
						return lbl.Layout(gtx)
					}
					if p.SuccessMsg != "" {
						lbl := material.Body2(th, p.SuccessMsg)
						lbl.Color = ui.SuccessColor
						return lbl.Layout(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btnText := "Update master password"
					if p.IsSaving {
						btnText = "Updating..."
					}
					return ui.PrimaryButton(gtx, th, &p.SaveMasterPasswordBtn, btnText)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
				layout.Rigid(material.Body1(th, "Appearance").Layout),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					prev := p.DarkMode.Value
					sw := material.Switch(th, &p.DarkMode, "Dark theme")
					d := sw.Layout(gtx)
					if darkMode != nil && prev != p.DarkMode.Value {
						*darkMode = p.DarkMode.Value
						changed = true
					}
					return d
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(th, "Credits")
					lbl.Color = th.Palette.ContrastBg
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(ui.MutedLabel(th, material.Body2(th, "Philopater Waheed | philopaterwaheed9@gmail.com")).Layout),
				layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return creditLinks(gtx, th, &p.GitHubBtn, &p.LinkedInBtn, &p.PortfolioBtn, &p.CVBtn)
				}),
			)
		})
	})

	return dims, changed, action
}
