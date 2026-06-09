package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

type LoginPage struct {
	EmailInput      widget.Editor
	PasswordInput   widget.Editor
	ShowPasswordBtn widget.Clickable
	ShowPassword    bool
	LoginBtn        widget.Clickable
	ForgotBtn       widget.Clickable
	BackBtn         widget.Clickable
	ErrorMsg        string
	SuccessMsg      string
	IsLoading       bool
}

func NewLoginPage() *LoginPage {
	return &LoginPage{
		EmailInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		PasswordInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
			Mask:       '*',
		},
	}
}

func (p *LoginPage) Reset() {
	p.EmailInput.SetText("")
	p.PasswordInput.SetText("")
	p.ShowPassword = false
	p.PasswordInput.Mask = '*'
	p.ErrorMsg = ""
	p.SuccessMsg = ""
	p.IsLoading = false
}
func (p *LoginPage) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(440)))
		children := []layout.FlexChild{
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.H4(th, "Log in")
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
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.GhostButton(gtx, th, &p.ForgotBtn, "Forgot password?")
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btnText := "Log in"
				if p.IsLoading {
					btnText = "Logging in..."
				}
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return ui.PrimaryButton(gtx, th, &p.LoginBtn, btnText)
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
