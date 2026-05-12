package ui

import (
	"image/color"
)

var BabyBlue = color.NRGBA{R: 173, G: 216, B: 230, A: 255}

func blend(a, b color.NRGBA, t uint8) color.NRGBA {
	u := uint16(255 - t)
	tt := uint16(t)
	return color.NRGBA{
		R: uint8((uint16(a.R)*u + uint16(b.R)*tt) / 255),
		G: uint8((uint16(a.G)*u + uint16(b.G)*tt) / 255),
		B: uint8((uint16(a.B)*u + uint16(b.B)*tt) / 255),
		A: uint8((uint16(a.A)*u + uint16(b.A)*tt) / 255),
	}
}

// for eye comfort
var BabyBlueBackground = blend(BabyBlue, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, 220)

