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
	List     widget.List

	MasterPasswordTabBtn  widget.Clickable
	AccountPasswordTabBtn widget.Clickable
	ActivePasswordTab     int

	CurrentMasterPasswordInput widget.Editor
	NewMasterPasswordInput     widget.Editor
	ConfirmMasterPasswordInput widget.Editor
	ShowCurrentMasterBtn       widget.Clickable
	ShowNewMasterBtn           widget.Clickable
	ShowConfirmMasterBtn       widget.Clickable
	ShowCurrentMaster          bool
	ShowNewMaster              bool
	ShowConfirmMaster          bool
	SaveMasterPasswordBtn      widget.Clickable

	CurrentAccountPasswordInput widget.Editor
	NewAccountPasswordInput     widget.Editor
	ConfirmAccountPasswordInput widget.Editor
	ShowCurrentAccountBtn       widget.Clickable
	ShowNewAccountBtn           widget.Clickable
	ShowConfirmAccountBtn       widget.Clickable
	ShowCurrentAccount          bool
	ShowNewAccount              bool
	ShowConfirmAccount          bool
	SaveAccountPasswordBtn      widget.Clickable

	GitHubBtn    widget.Clickable
	LinkedInBtn  widget.Clickable
	PortfolioBtn widget.Clickable
	CVBtn        widget.Clickable
	ErrorMsg     string
	SuccessMsg   string
	IsSaving     bool

	initialized bool
}

type SettingsAction struct {
	UpdateMasterPassword  bool
	CurrentMasterPassword string
	NewMasterPassword     string
	ConfirmMasterPassword string

	UpdateAccountPassword  bool
	CurrentAccountPassword string
	NewAccountPassword     string
	ConfirmAccountPassword string
}

const (
	settingsMasterPasswordTab = iota
	settingsAccountPasswordTab
)

func NewSettingsPage() *SettingsPage {
	return &SettingsPage{
		List: widget.List{List: layout.List{Axis: layout.Vertical}},
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
		CurrentAccountPasswordInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
			Mask:       '*',
		},
		NewAccountPasswordInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
			Mask:       '*',
		},
		ConfirmAccountPasswordInput: widget.Editor{
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
	p.ShowCurrentMaster = false
	p.ShowNewMaster = false
	p.ShowConfirmMaster = false
	p.CurrentMasterPasswordInput.Mask = '*'
	p.NewMasterPasswordInput.Mask = '*'
	p.ConfirmMasterPasswordInput.Mask = '*'
}

func (p *SettingsPage) ClearAccountPasswordFields() {
	p.CurrentAccountPasswordInput.SetText("")
	p.NewAccountPasswordInput.SetText("")
	p.ConfirmAccountPasswordInput.SetText("")
	p.ShowCurrentAccount = false
	p.ShowNewAccount = false
	p.ShowConfirmAccount = false
	p.CurrentAccountPasswordInput.Mask = '*'
	p.NewAccountPasswordInput.Mask = '*'
	p.ConfirmAccountPasswordInput.Mask = '*'
}

func (p *SettingsPage) Reset() {
	p.ClearMasterPasswordFields()
	p.ClearAccountPasswordFields()
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

	for p.MasterPasswordTabBtn.Clicked(gtx) {
		p.ActivePasswordTab = settingsMasterPasswordTab
		p.ErrorMsg = ""
		p.SuccessMsg = ""
	}
	for p.AccountPasswordTabBtn.Clicked(gtx) {
		p.ActivePasswordTab = settingsAccountPasswordTab
		p.ErrorMsg = ""
		p.SuccessMsg = ""
	}
	if p.SaveMasterPasswordBtn.Clicked(gtx) {
		action.UpdateMasterPassword = true
		action.CurrentMasterPassword = p.CurrentMasterPasswordInput.Text()
		action.NewMasterPassword = p.NewMasterPasswordInput.Text()
		action.ConfirmMasterPassword = p.ConfirmMasterPasswordInput.Text()
	}
	if p.SaveAccountPasswordBtn.Clicked(gtx) {
		action.UpdateAccountPassword = true
		action.CurrentAccountPassword = p.CurrentAccountPasswordInput.Text()
		action.NewAccountPassword = p.NewAccountPasswordInput.Text()
		action.ConfirmAccountPassword = p.ConfirmAccountPasswordInput.Text()
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
			return p.List.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.H6(th, "Account")
						lbl.Color = th.Palette.ContrastBg
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
					layout.Rigid(ui.MutedLabel(th, material.Body2(th, "Signed in as: "+email)).Layout),
					layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return passwordTab(gtx, th, &p.MasterPasswordTabBtn, "Master password", p.ActivePasswordTab == settingsMasterPasswordTab)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return passwordTab(gtx, th, &p.AccountPasswordTabBtn, "Account password", p.ActivePasswordTab == settingsAccountPasswordTab)
							}),
						)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if p.ActivePasswordTab == settingsAccountPasswordTab {
							return p.accountPasswordFields(gtx, th)
						}
						return p.masterPasswordFields(gtx, th)
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
	})

	return dims, changed, action
}

func passwordTab(gtx layout.Context, th *material.Theme, click *widget.Clickable, label string, active bool) layout.Dimensions {
	if active {
		return ui.PrimaryButton(gtx, th, click, label)
	}
	return ui.SecondaryButton(gtx, th, click, label)
}

func (p *SettingsPage) masterPasswordFields(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.PasswordEditor(gtx, th, &p.CurrentMasterPasswordInput, "Current master password", &p.ShowCurrentMasterBtn, &p.ShowCurrentMaster)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.PasswordEditor(gtx, th, &p.NewMasterPasswordInput, "New master password", &p.ShowNewMasterBtn, &p.ShowNewMaster)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.PasswordEditor(gtx, th, &p.ConfirmMasterPasswordInput, "Confirm new master password", &p.ShowConfirmMasterBtn, &p.ShowConfirmMaster)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.passwordMessage(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnText := "Update master password"
			if p.IsSaving {
				btnText = "Updating..."
			}
			return ui.PrimaryButton(gtx, th, &p.SaveMasterPasswordBtn, btnText)
		}),
	)
}

func (p *SettingsPage) accountPasswordFields(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.PasswordEditor(gtx, th, &p.CurrentAccountPasswordInput, "Current account password", &p.ShowCurrentAccountBtn, &p.ShowCurrentAccount)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.PasswordEditor(gtx, th, &p.NewAccountPasswordInput, "New account password", &p.ShowNewAccountBtn, &p.ShowNewAccount)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.PasswordEditor(gtx, th, &p.ConfirmAccountPasswordInput, "Confirm new account password", &p.ShowConfirmAccountBtn, &p.ShowConfirmAccount)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.passwordMessage(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnText := "Update account password"
			if p.IsSaving {
				btnText = "Updating..."
			}
			return ui.PrimaryButton(gtx, th, &p.SaveAccountPasswordBtn, btnText)
		}),
	)
}

func (p *SettingsPage) passwordMessage(gtx layout.Context, th *material.Theme) layout.Dimensions {
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
}
