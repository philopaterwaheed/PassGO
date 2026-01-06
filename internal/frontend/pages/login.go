package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type LoginPage struct {
	EmailInput    widget.Editor
	PasswordInput widget.Editor
	LoginBtn      widget.Clickable
	BackBtn       widget.Clickable
	ErrorMsg      string
	SuccessMsg    string
	IsLoading     bool
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
	p.ErrorMsg = ""
	p.SuccessMsg = ""
	p.IsLoading = false
}
func (p *LoginPage) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		children := []layout.FlexChild{
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.H4(th, "Login").Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		}

		// Show error message if any
		if p.ErrorMsg != "" {
			children = append(children,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					l := material.Body2(th, p.ErrorMsg)
					l.Color = th.Palette.ContrastBg // Red color
					return l.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			)
		}

		// Show success message if any
		if p.SuccessMsg != "" {
			children = append(children,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					l := material.Body2(th, p.SuccessMsg)
					l.Color = th.Palette.ContrastFg // Green-ish color
					return l.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			)
		}

		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(th, &p.EmailInput, "Email")
				return e.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				e := material.Editor(th, &p.PasswordInput, "Password")
				return e.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btnText := "Login"
				if p.IsLoading {
					btnText = "Loading..."
				}
				btn := material.Button(th, &p.LoginBtn, btnText)
				if p.IsLoading {
					btn.Background = th.Palette.Bg
				}
				return btn.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Button(th, &p.BackBtn, "Back").Layout(gtx)
			}),
		)

		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx, children...)
	})
}
