package ui

import (
	"crypto/rand"
	"image"
	"image/color"
	"math/big"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

const generatedPasswordLength = 20

var generatedPasswordSets = []string{
	"abcdefghijkmnopqrstuvwxyz",
	"ABCDEFGHJKLMNPQRSTUVWXYZ",
	"23456789",
	"!@#$%&*?-_",
}

func GeneratePassword() (string, error) {
	all := ""
	for _, set := range generatedPasswordSets {
		all += set
	}

	password := make([]byte, 0, generatedPasswordLength)
	for _, set := range generatedPasswordSets {
		char, err := randomChar(set)
		if err != nil {
			return "", err
		}
		password = append(password, char)
	}

	for len(password) < generatedPasswordLength {
		char, err := randomChar(all)
		if err != nil {
			return "", err
		}
		password = append(password, char)
	}

	for i := len(password) - 1; i > 0; i-- {
		j, err := randomIndex(i + 1)
		if err != nil {
			return "", err
		}
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

func randomChar(chars string) (byte, error) {
	idx, err := randomIndex(len(chars))
	if err != nil {
		return 0, err
	}
	return chars[idx], nil
}

func randomIndex(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func PasswordEditor(gtx layout.Context, th *material.Theme, editor *widget.Editor, hint string, toggle *widget.Clickable, shown *bool) layout.Dimensions {
	return passwordEditor(gtx, th, editor, hint, toggle, shown, nil, false)
}

func PasswordEditorWithGenerate(gtx layout.Context, th *material.Theme, editor *widget.Editor, hint string, toggle *widget.Clickable, shown *bool, generate *widget.Clickable) layout.Dimensions {
	return passwordEditor(gtx, th, editor, hint, toggle, shown, generate, false)
}

func passwordEditor(gtx layout.Context, th *material.Theme, editor *widget.Editor, hint string, toggle *widget.Clickable, shown *bool, generate *widget.Clickable, fullWidth bool) layout.Dimensions {
	for toggle.Clicked(gtx) {
		*shown = !*shown
	}
	if *shown {
		editor.Mask = 0
	} else {
		editor.Mask = '*'
	}

	editorChild := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		if !fullWidth {
			buttonCount := 1
			if generate != nil {
				buttonCount++
			}
			gtx.Constraints.Max.X = maxInt(0, gtx.Constraints.Max.X-(gtx.Dp(unit.Dp(34))*buttonCount))
			gtx.Constraints.Min.X = minInt(gtx.Constraints.Min.X, gtx.Constraints.Max.X)
		}
		e := material.Editor(th, editor, hint)
		e.TextSize = unit.Sp(16)
		return e.Layout(gtx)
	})
	if fullWidth {
		editorChild = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			e := material.Editor(th, editor, hint)
			e.TextSize = unit.Sp(16)
			return e.Layout(gtx)
		})
	}

	children := []layout.FlexChild{
		editorChild,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return EyeButton(gtx, th, toggle, *shown)
		}),
	}
	if generate != nil {
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return GeneratePasswordButton(gtx, th, generate)
		}))
	}

	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
}

func EyeButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, shown bool) layout.Dimensions {
	label := "S"
	size := unit.Sp(13)
	if shown {
		label = "H"
		size = unit.Sp(14)
	}
	return SmallPasswordButton(gtx, th, click, label, size)
}

func GeneratePasswordButton(gtx layout.Context, th *material.Theme, click *widget.Clickable) layout.Dimensions {
	return SmallPasswordButton(gtx, th, click, "G", unit.Sp(14))
}

func SmallPasswordButton(gtx layout.Context, th *material.Theme, click *widget.Clickable, text string, textSize unit.Sp) layout.Dimensions {
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
			label := material.Label(th, textSize, text)
			label.Color = iconColor(th.Palette.Fg)
			return label.Layout(gtx)
		})
	})
}

func iconColor(col color.NRGBA) color.NRGBA {
	col.A = 190
	return col
}
