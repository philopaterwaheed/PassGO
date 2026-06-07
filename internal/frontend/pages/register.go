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
	p.ErrorMsg = ""
	p.SuccessMsg = ""
	p.IsLoading = false
}

func (p *RegisterPage) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(420)))
		children := []layout.FlexChild{
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, material.H4(th, "Create account").Layout)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		}

		// Show error message if any
		if p.ErrorMsg != "" {
			children = append(children,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, material.Body2(th, p.ErrorMsg).Layout)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			)
		}

		// Show success message if any
		if p.SuccessMsg != "" {
			children = append(children,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, material.Body2(th, p.SuccessMsg).Layout)
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
				e := material.Editor(th, &p.PasswordInput, "Password")
				e.TextSize = unit.Sp(16)
				return e.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(th, &p.ConfirmPasswordInput, "Confirm password")
				e.TextSize = unit.Sp(16)
				return e.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(th, &p.MasterPasswordInput, "Master password")
				e.TextSize = unit.Sp(16)
				return e.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(th, &p.ConfirmMasterPasswordInput, "Confirm master password")
				e.TextSize = unit.Sp(16)
				return e.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btnText := "Create account"
				if p.IsLoading {
					btnText = "Creating..."
				}
				btn := material.Button(th, &p.RegisterBtn, btnText)
				btn.TextSize = unit.Sp(16)
				return btn.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.OutlinedButton(gtx, th, &p.BackBtn, "Back", th.Palette.Fg, th.Palette.ContrastBg)
			}),
		)

		return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx, children...)
		})
	})
}
