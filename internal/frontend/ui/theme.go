package ui

import (
	"image/color"

	"gioui.org/widget/material"
)

var darkForeground = color.NRGBA{R: 236, G: 246, B: 255, A: 255}
var contrastText = color.NRGBA{R: 10, G: 31, B: 38, A: 255}

func ApplyTheme(th *material.Theme, dark bool) {
	if th == nil {
		return
	}

	if dark {
		th.Palette.Bg = color.NRGBA{R: 18, G: 30, B: 36, A: 255}
		th.Palette.Fg = darkForeground
		th.Palette.ContrastBg = BabyBlue
		th.Palette.ContrastFg = contrastText
		SurfaceColor = color.NRGBA{R: 27, G: 43, B: 50, A: 255}
		SurfaceAltColor = color.NRGBA{R: 34, G: 53, B: 61, A: 255}
		BorderColor = color.NRGBA{R: 78, G: 114, B: 126, A: 255}
		MutedColor = color.NRGBA{R: 180, G: 205, B: 214, A: 255}
		DangerColor = color.NRGBA{R: 238, G: 111, B: 111, A: 255}
		SuccessColor = color.NRGBA{R: 116, G: 220, B: 173, A: 255}
		return
	}

	th.Palette.Bg = BabyBlueBackground
	th.Palette.Fg = contrastText
	th.Palette.ContrastBg = BabyBlue
	th.Palette.ContrastFg = contrastText
	SurfaceColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	SurfaceAltColor = color.NRGBA{R: 232, G: 245, B: 249, A: 255}
	BorderColor = color.NRGBA{R: 180, G: 218, B: 227, A: 255}
	MutedColor = color.NRGBA{R: 79, G: 107, B: 116, A: 255}
	DangerColor = color.NRGBA{R: 183, G: 48, B: 58, A: 255}
	SuccessColor = color.NRGBA{R: 29, G: 132, B: 83, A: 255}
}
