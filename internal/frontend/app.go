package frontend

import (
	"image/color"
	"log"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/philopaterwaheed/passGO/internal/frontend/api"
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

	// Initialize API client
	// Default to localhost:8080, can be configured
	apiClient := api.NewClient("https://fantastic-halibut-756rjg76p7g2q9p-8080.app.github.dev")

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

			// Handle login
			if loginPage.LoginBtn.Clicked(gtx) && !loginPage.IsLoading {
				email := loginPage.EmailInput.Text()
				password := loginPage.PasswordInput.Text()

				if email == "" || password == "" {
					loginPage.ErrorMsg = "Email and password are required"
				} else {
					loginPage.IsLoading = true
					loginPage.ErrorMsg = ""
					loginPage.SuccessMsg = ""

					// Call backend API in goroutine
					go func() {
						resp, err := apiClient.Login(email, password)
						if err != nil {
							loginPage.ErrorMsg = err.Error()
							loginPage.IsLoading = false
							w.Invalidate()
							return
						}

						loginPage.SuccessMsg = "Login successful! Welcome, " + resp.User.Email
						loginPage.IsLoading = false
						// TODO: Navigate to main app and store token
						log.Printf("Logged in successfully: %+v", resp.User)
						w.Invalidate()
					}()
				}
			}

			// Handle registration
			if registerPage.RegisterBtn.Clicked(gtx) && !registerPage.IsLoading {
				email := registerPage.EmailInput.Text()
				password := registerPage.PasswordInput.Text()
				confirmPassword := registerPage.ConfirmPasswordInput.Text()

				if email == "" || password == "" || confirmPassword == "" {
					registerPage.ErrorMsg = "All fields are required"
				} else if !strings.Contains(email, "@") {
					registerPage.ErrorMsg = "Invalid email address"
				} else if len(password) < 8 {
					registerPage.ErrorMsg = "Password must be at least 8 characters"
				} else if password != confirmPassword {
					registerPage.ErrorMsg = "Passwords do not match"
				} else {
					registerPage.IsLoading = true
					registerPage.ErrorMsg = ""
					registerPage.SuccessMsg = ""

					// Call backend API in goroutine
					go func() {
						resp, err := apiClient.Signup(email, password)
						if err != nil {
							registerPage.ErrorMsg = err.Error()
							registerPage.IsLoading = false
							w.Invalidate()
							return
						}

						registerPage.SuccessMsg = resp.Message
						if registerPage.SuccessMsg == "" {
							registerPage.SuccessMsg = "Registration successful! Please check your email."
						}
						registerPage.IsLoading = false
						// Clear password fields after successful registration
						registerPage.PasswordInput.SetText("")
						registerPage.ConfirmPasswordInput.SetText("")
						log.Printf("Registered successfully: %+v", resp.User)
						w.Invalidate()
					}()
				}
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
