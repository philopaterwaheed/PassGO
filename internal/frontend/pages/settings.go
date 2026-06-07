package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type SettingsPage struct {
	DarkMode widget.Bool

	CurrentMasterPasswordInput widget.Editor
	NewMasterPasswordInput     widget.Editor
	ConfirmMasterPasswordInput widget.Editor
	SaveMasterPasswordBtn      widget.Clickable
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

	dims := layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(material.Body1(th, "Account").Layout),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(material.Body2(th, "Signed in as: "+email).Layout),
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
					return material.Body2(th, p.ErrorMsg).Layout(gtx)
				}
				if p.SuccessMsg != "" {
					return material.Body2(th, p.SuccessMsg).Layout(gtx)
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btnText := "Update master password"
				if p.IsSaving {
					btnText = "Updating..."
				}
				btn := material.Button(th, &p.SaveMasterPasswordBtn, btnText)
				btn.TextSize = unit.Sp(14)
				return btn.Layout(gtx)
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
		)
	})

	return dims, changed, action
}
