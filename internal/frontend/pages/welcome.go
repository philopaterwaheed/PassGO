package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

type WelcomePage struct {
	LoginBtn    widget.Clickable
	RegisterBtn widget.Clickable
}

func NewWelcomePage() *WelcomePage {
	return &WelcomePage{}
}

func (p *WelcomePage) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(420)))
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(material.H4(th, "PassGO").Layout),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(material.Body1(th, "Minimal, fast password manager").Layout),
				layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &p.LoginBtn, "Log in")
					btn.Background = th.Palette.ContrastBg
					btn.TextSize = unit.Sp(16)
					return btn.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx2 := gtx
					return ui.OutlinedButton(gtx2, th, &p.RegisterBtn, "Create account", th.Palette.Fg, th.Palette.ContrastBg)
				}),
			)
		})
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
