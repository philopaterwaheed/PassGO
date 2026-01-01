package frontend

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	. "github.com/philopaterwaheed/passGO/internal/frontend/pages"
)

// Run starts the Gio desktop/web application
func Run() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("PassGO - Password Manager"),
			app.Size(unit.Dp(800), unit.Dp(600)),
		)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	var loginBtn widget.Clickable
	var registerBtn widget.Clickable

	loginPage := NewLoginPage()
	registerPage := NewRegisterPage()
	currentPage := "welcome"

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if loginBtn.Clicked(gtx) {
				currentPage = "login"
			}
			if registerBtn.Clicked(gtx) {
				currentPage = "register"
			}
			if loginPage.BackBtn.Clicked(gtx) {
				currentPage = "welcome"
				loginPage.Reset()
			}
			if registerPage.BackBtn.Clicked(gtx) {
				currentPage = "welcome"
				registerPage.Reset()
			}

			if currentPage == "login" {
				loginPage.Layout(gtx, th)
			} else if currentPage == "register" {
				registerPage.Layout(gtx, th)
			} else {
				layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Vertical,
						Spacing:   layout.SpaceBetween,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							l := material.H3(th, "Welcome to PassGO")
							l.Alignment = text.Middle
							l.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
							return l.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							l := material.Body1(th, "Your secure password manager")
							l.Alignment = text.Middle
							l.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
							return l.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:    layout.Horizontal,
								Spacing: layout.SpaceEvenly,
							}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Button(th, &loginBtn, "Login").Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Button(th, &registerBtn, "Register").Layout(gtx)
								}),
							)
						}),
					)
				})
			}

			e.Frame(gtx.Ops)
		}
	}
}
