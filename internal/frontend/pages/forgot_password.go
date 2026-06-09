package pages

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

type ForgotPasswordPage struct {
	EmailInput widget.Editor
	SendBtn    widget.Clickable
	BackBtn    widget.Clickable
	ErrorMsg   string
	SuccessMsg string
	IsLoading  bool
}

func NewForgotPasswordPage() *ForgotPasswordPage {
	return &ForgotPasswordPage{
		EmailInput: widget.Editor{SingleLine: true, Submit: true},
	}
}

func (p *ForgotPasswordPage) Reset() {
	p.EmailInput.SetText("")
	p.ErrorMsg = ""
	p.SuccessMsg = ""
	p.IsLoading = false
}

func (p *ForgotPasswordPage) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(440)))
		children := []layout.FlexChild{
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.H4(th, "Reset password")
				lbl.Color = th.Palette.ContrastBg
				return layout.Center.Layout(gtx, lbl.Layout)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ui.MutedLabel(th, material.Body2(th, "We will email you a reset link.")).Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
		}

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
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				e := material.Editor(th, &p.EmailInput, "Email")
				e.TextSize = unit.Sp(16)
				return e.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btnText := "Send reset link"
				if p.IsLoading {
					btnText = "Sending..."
				}
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return ui.PrimaryButton(gtx, th, &p.SendBtn, btnText)
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
