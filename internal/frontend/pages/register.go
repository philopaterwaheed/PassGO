package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

type RegisterPage struct {
	EmailInput                 widget.Editor
	PasswordInput              widget.Editor
	ConfirmPasswordInput       widget.Editor
	MasterPasswordInput        widget.Editor
	ConfirmMasterPasswordInput widget.Editor
	ShowPasswordBtn            widget.Clickable
	ShowConfirmPasswordBtn     widget.Clickable
	ShowMasterPasswordBtn      widget.Clickable
	ShowConfirmMasterBtn       widget.Clickable
	ShowPassword               bool
	ShowConfirmPassword        bool
	ShowMasterPassword         bool
	ShowConfirmMasterPassword  bool
	RegisterBtn                widget.Clickable
	BackBtn                    widget.Clickable
	ErrorMsg                   string
	SuccessMsg                 string
	IsLoading                  bool
}

func NewRegisterPage() *RegisterPage {
	return &RegisterPage{
		EmailInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		PasswordInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
			Mask:       '*',
		},
		ConfirmPasswordInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
			Mask:       '*',
		},
		MasterPasswordInput: widget.Editor{
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

func (p *RegisterPage) Reset() {
	p.EmailInput.SetText("")
	p.PasswordInput.SetText("")
	p.ConfirmPasswordInput.SetText("")
	p.MasterPasswordInput.SetText("")
	p.ConfirmMasterPasswordInput.SetText("")
	p.ShowPassword = false
	p.ShowConfirmPassword = false
	p.ShowMasterPassword = false
	p.ShowConfirmMasterPassword = false
	p.PasswordInput.Mask = '*'
	p.ConfirmPasswordInput.Mask = '*'
	p.MasterPasswordInput.Mask = '*'
	p.ConfirmMasterPasswordInput.Mask = '*'
	p.ErrorMsg = ""
	p.SuccessMsg = ""
	p.IsLoading = false
}

func (p *RegisterPage) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(460)))
		children := []layout.FlexChild{
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.H4(th, "Create account")
				lbl.Color = th.Palette.ContrastBg
				return layout.Center.Layout(gtx, lbl.Layout)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		}

		// Show error message if any
		if p.ErrorMsg != "" {
			children = append(children,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th, p.ErrorMsg)
					lbl.Color = ui.DangerColor
					return layout.Center.Layout(gtx, lbl.Layout)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			)
		}

		// Show success message if any
		if p.SuccessMsg != "" {
			children = append(children,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th, p.SuccessMsg)
					lbl.Color = ui.SuccessColor
					return layout.Center.Layout(gtx, lbl.Layout)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			)
		}

		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(th, &p.EmailInput, "Email")
				e.TextSize = unit.Sp(16)
				return e.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.PasswordEditor(gtx, th, &p.PasswordInput, "Password", &p.ShowPasswordBtn, &p.ShowPassword)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.PasswordEditor(gtx, th, &p.ConfirmPasswordInput, "Confirm password", &p.ShowConfirmPasswordBtn, &p.ShowConfirmPassword)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.PasswordEditor(gtx, th, &p.MasterPasswordInput, "Master password", &p.ShowMasterPasswordBtn, &p.ShowMasterPassword)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.PasswordEditor(gtx, th, &p.ConfirmMasterPasswordInput, "Confirm master password", &p.ShowConfirmMasterBtn, &p.ShowConfirmMasterPassword)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btnText := "Create account"
				if p.IsLoading {
					btnText = "Creating..."
				}
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return ui.PrimaryButton(gtx, th, &p.RegisterBtn, btnText)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return ui.SecondaryButton(gtx, th, &p.BackBtn, "Back")
			}),
		)

		return ui.SurfacePanel(gtx, ui.SurfaceColor, ui.RadiusLarge, layout.UniformInset(unit.Dp(24)), func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx, children...)
		})
	})
}
