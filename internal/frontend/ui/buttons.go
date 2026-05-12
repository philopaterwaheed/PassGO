package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

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
	btn.CornerRadius = unit.Dp(10)

	border := widget.Border{
		Color:        borderCol,
		CornerRadius: btn.CornerRadius,
		Width:        unit.Dp(1),
	}

	return border.Layout(gtx, btn.Layout)
}