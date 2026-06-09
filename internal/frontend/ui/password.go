package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func PasswordEditor(gtx layout.Context, th *material.Theme, editor *widget.Editor, hint string, toggle *widget.Clickable, shown *bool) layout.Dimensions {
	for toggle.Clicked(gtx) {
		*shown = !*shown
	}
	if *shown {
		editor.Mask = 0
	} else {
		editor.Mask = '*'
	}

	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			e := material.Editor(th, editor, hint)
			e.TextSize = unit.Sp(16)
			return e.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return EyeButton(gtx, th, toggle, *shown)
		}),
	)
}

func EyeButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, shown bool) layout.Dimensions {
	size := gtx.Dp(unit.Dp(34))
	gtx.Constraints.Min.X = size
	gtx.Constraints.Max.X = size
	gtx.Constraints.Min.Y = size
	gtx.Constraints.Max.Y = size

	return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if click.Hovered() || click.Pressed() {
			rect := image.Rectangle{Max: image.Pt(size, size)}
			paint.FillShape(gtx.Ops, SurfaceAltColor, clip.UniformRRect(rect, gtx.Dp(RadiusSmall)).Op(gtx.Ops))
		}

		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			label := material.Label(th, unit.Sp(13), "S")
			if shown {
				label = material.Label(th, unit.Sp(14), "H")
			}
			label.Color = iconColor(th.Palette.Fg)
			return label.Layout(gtx)
		})
	})
}

func iconColor(col color.NRGBA) color.NRGBA {
	col.A = 190
	return col
}
