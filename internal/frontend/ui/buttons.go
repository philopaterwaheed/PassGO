package ui

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

const (
	RadiusSmall unit.Dp = 8
	Radius      unit.Dp = 12
	RadiusLarge unit.Dp = 16
)

func PageInset(gtx layout.Context, w layout.Widget) layout.Dimensions {
	pad := unit.Dp(16)
	if gtx.Constraints.Max.X > gtx.Dp(unit.Dp(800)) {
		pad = unit.Dp(24)
	}
	return layout.UniformInset(pad).Layout(gtx, w)
}

func ConstrainWidth(gtx layout.Context, max unit.Dp, w layout.Widget) layout.Dimensions {
	gtx.Constraints.Max.X = minInt(gtx.Constraints.Max.X, gtx.Dp(max))
	return w(gtx)
}

func SurfacePanel(gtx layout.Context, bg color.NRGBA, radius unit.Dp, inset layout.Inset, w layout.Widget) layout.Dimensions {
	content := func(gtx layout.Context) layout.Dimensions {
		return inset.Layout(gtx, w)
	}
	return surface(gtx, bg, BorderColor, radius, unit.Dp(1), content)
}

func Surface(gtx layout.Context, bg color.NRGBA, radius unit.Dp, w layout.Widget) layout.Dimensions {
	return surface(gtx, bg, BorderColor, radius, unit.Dp(1), w)
}

func surface(gtx layout.Context, bg, border color.NRGBA, radius unit.Dp, width unit.Dp, w layout.Widget) layout.Dimensions {
	mac := op.Record(gtx.Ops)
	dims := w(gtx)
	call := mac.Stop()

	bounds := image.Rectangle{Max: dims.Size}
	defer clip.UniformRRect(bounds, gtx.Dp(radius)).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, bg)
	widget.Border{
		Color:        border,
		CornerRadius: radius,
		Width:        width,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: dims.Size}
	})
	call.Add(gtx.Ops)
	return dims
}

func PrimaryButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, label string) layout.Dimensions {
	btn := material.Button(th, click, label)
	btn.TextSize = unit.Sp(15)
	btn.CornerRadius = RadiusSmall
	btn.Background = th.Palette.ContrastBg
	btn.Color = th.Palette.ContrastFg
	btn.Inset = layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(14), Right: unit.Dp(14)}
	return btn.Layout(gtx)
}

func SecondaryButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, label string) layout.Dimensions {
	return OutlinedButton(gtx, th, click, label, th.Palette.Fg, BorderColor)
}

func GhostButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, label string) layout.Dimensions {
	btn := material.Button(th, click, label)
	btn.Background = color.NRGBA{}
	btn.Color = th.Palette.Fg
	btn.CornerRadius = RadiusSmall
	btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(10), Right: unit.Dp(10)}
	btn.TextSize = unit.Sp(14)
	return btn.Layout(gtx)
}

func CopyButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, copied bool) layout.Dimensions {
	label := "Copy"
	col := th.Palette.Fg
	borderCol := BorderColor
	if copied {
		label = "Copied"
		col = SuccessColor
		borderCol = SuccessColor
	}

	btn := material.Button(th, click, label)
	btn.Background = color.NRGBA{}
	btn.Color = col
	btn.CornerRadius = RadiusSmall
	btn.TextSize = unit.Sp(13)
	btn.Inset = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(10), Right: unit.Dp(10)}

	border := widget.Border{
		Color:        borderCol,
		CornerRadius: btn.CornerRadius,
		Width:        unit.Dp(1),
	}

	return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = maxInt(gtx.Constraints.Min.X, gtx.Dp(unit.Dp(72)))
		gtx.Constraints.Min.Y = maxInt(gtx.Constraints.Min.Y, gtx.Dp(unit.Dp(36)))
		return btn.Layout(gtx)
	})
}

func LinkButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, icon string, label string) layout.Dimensions {
	return linkButton(gtx, th, click, icon, label, unit.Dp(78), layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(8), Right: unit.Dp(8)})
}

func CompactLinkButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, icon string, label string) layout.Dimensions {
	return linkButton(gtx, th, click, icon, label, unit.Dp(102), layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(8), Right: unit.Dp(8)})
}

func linkButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, icon string, label string, minWidth unit.Dp, inset layout.Inset) layout.Dimensions {
	bg := color.NRGBA{}
	if click.Hovered() || click.Pressed() {
		bg = SurfaceAltColor
	}

	return surface(gtx, bg, BorderColor, RadiusSmall, unit.Dp(1), func(gtx layout.Context) layout.Dimensions {
		return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = maxInt(gtx.Constraints.Min.Y, gtx.Dp(unit.Dp(34)))
			if uiWidth := gtx.Dp(minWidth); gtx.Constraints.Max.X >= uiWidth {
				gtx.Constraints.Min.X = maxInt(gtx.Constraints.Min.X, uiWidth)
			}
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Caption(th, icon)
							lbl.Color = th.Palette.ContrastBg
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Caption(th, label)
							lbl.Color = th.Palette.Fg
							return lbl.Layout(gtx)
						}),
					)
				})
			})
		})
	})
}

func MutedLabel(th *material.Theme, label material.LabelStyle) material.LabelStyle {
	label.Color = MutedColor
	return label
}

// OutlinedButton draws an outline-style button with a custom color
func OutlinedButton(
	gtx layout.Context,
	th *material.Theme,
	click *widget.Clickable,
	label string,
	col color.NRGBA,
	borderCol color.NRGBA,
) layout.Dimensions {

	btn := material.Button(th, click, label)
	btn.Background = color.NRGBA{} // transparent
	btn.Color = col
	btn.CornerRadius = RadiusSmall
	btn.TextSize = unit.Sp(15)
	btn.Inset = layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(14), Right: unit.Dp(14)}

	border := widget.Border{
		Color:        borderCol,
		CornerRadius: btn.CornerRadius,
		Width:        unit.Dp(1),
	}

	return border.Layout(gtx, btn.Layout)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
