package pages

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"image"
)

func SecurityInfoCard(gtx layout.Context, th *material.Theme) layout.Dimensions {
	radius := gtx.Dp(unit.Dp(8))
	bg := th.Palette.ContrastBg
	bg.A = 30 // Soft transparent background

	return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		content := func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body1(th, "🔒")
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(th, "All values are safely encrypted on your device using AES-256-GCM. Your master password is never stored.")
						lbl.Color = th.Palette.Fg
						lbl.Color.A = 180
						return lbl.Layout(gtx)
					}),
				)
			})
		}

		mac := op.Record(gtx.Ops)
		dims := content(gtx)
		call := mac.Stop()

		defer clip.UniformRRect(image.Rectangle{Max: dims.Size}, radius).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, bg)
		call.Add(gtx.Ops)
		return dims
	})
}
