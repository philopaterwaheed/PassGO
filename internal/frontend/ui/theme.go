package ui

import (
	"image/color"

	"gioui.org/widget/material"
)

var babyBlueDarkBackground = blend(BabyBlue, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, 220)
var darkForeground = color.NRGBA{R: 236, G: 246, B: 255, A: 255}
var contrastText = color.NRGBA{R: 0, G: 0, B: 0, A: 255}


func ApplyTheme(th *material.Theme, dark bool) {
	if th == nil {
		return
	}

	if dark {
		th.Palette.Bg = babyBlueDarkBackground
		th.Palette.Fg = darkForeground
		th.Palette.ContrastBg = BabyBlue
		th.Palette.ContrastFg = contrastText
		return
	}

	th.Palette.Bg = BabyBlueBackground
	th.Palette.Fg = contrastText
	th.Palette.ContrastBg = BabyBlue
	th.Palette.ContrastFg = contrastText
}
