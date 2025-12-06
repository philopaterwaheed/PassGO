package frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Run starts the Fyne desktop application
func Run() {
	myApp := app.New()
	myWindow := myApp.NewWindow("PassGO - Password Manager")

	// Create the main UI
	content := createMainUI()
	myWindow.SetContent(content)

	// Set window size
	myWindow.Resize(fyne.NewSize(800, 600))

	// Show and run
	myWindow.ShowAndRun()
}

// createMainUI creates the main user interface
func createMainUI() fyne.CanvasObject {
	// Welcome label
	welcomeLabel := widget.NewLabel("Welcome to PassGO")
	welcomeLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Description
	descLabel := widget.NewLabel("Your secure password manager")

	// Placeholder buttons
	loginBtn := widget.NewButton("Login", func() {
		// TODO: Implement login functionality
	})

	registerBtn := widget.NewButton("Register", func() {
		// TODO: Implement register functionality
	})

	// Layout
	content := container.NewVBox(
		welcomeLabel,
		descLabel,
		container.NewHBox(loginBtn, registerBtn),
	)

	return container.NewCenter(content)
}
