package frontend

import (
	"log"
	"os"
	"runtime"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/joho/godotenv"
	"github.com/philopaterwaheed/passGO/internal/frontend/api"
	"github.com/philopaterwaheed/passGO/internal/frontend/pages"
	"github.com/philopaterwaheed/passGO/internal/frontend/state"
	"github.com/philopaterwaheed/passGO/internal/frontend/storage"
	"github.com/philopaterwaheed/passGO/internal/frontend/ui"
)

// Run starts the Gio desktop/web application
func Run() {
	log.Printf("passGo starting on %s/%s", runtime.GOOS, runtime.GOARCH)
	go func() {
		w := new(app.Window)
		// Desktop gets a sensible default size; mobile/web will size the surface.
		if runtime.GOOS != "android" && runtime.GOOS != "ios" && runtime.GOOS != "js" {
			w.Option(
				app.Title("PassGO - Password Manager"),
				app.Size(unit.Dp(1024), unit.Dp(720)),
			)
		} else {
			w.Option(app.Title("PassGO - Password Manager"))
		}
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		log.Printf("passGo stopped")
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	log.Printf("passGo initializing")
	// Best-effort: make local desktop/mobile runs pick up .env without requiring
	// manual exporting. (On WASM it will typically just fail and be ignored.)
	_ = godotenv.Load()

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	st := &state.AppState{Route: state.RouteWelcome, Nav: state.NavVault, DarkMode: true}
	ui.ApplyTheme(th, st.DarkMode)

	welcomePage := pages.NewWelcomePage()
	loginPage := pages.NewLoginPage()
	registerPage := pages.NewRegisterPage()
	forgotPasswordPage := pages.NewForgotPasswordPage()
	vaultListPage := pages.NewVaultListPage()
	vaultAddPage := pages.NewVaultAddPage()
	vaultDetailPage := pages.NewVaultDetailPage()
	settingsPage := pages.NewSettingsPage()
	shell := ui.NewShell()

	apiBaseURL := strings.TrimRight(os.Getenv("PASSGO_API_BASE_URL"), "/")
	if apiBaseURL == "" {
		// apiBaseURL = "https://passgo.leapcell.app"
		apiBaseURL = "https://curly-memory-xp79gjr7q5gfpr46-8080.app.github.dev"
	}
	apiClient := api.NewClient(apiBaseURL)
	sessionStore := storage.NewSessionStore()

	if sess, err := sessionStore.Load(); err == nil && sess.Token != "" {
		log.Printf("PassGO restored saved session")
		st.Auth.Token = sess.Token
		st.Auth.Email = sess.Email
		st.Auth.MasterPassword = ""
		apiClient.Token = sess.Token
		st.Nav = state.NavVault
		st.Route = state.RouteVaultList

		// Best-effort validation/refresh of the restored session.
		go func() {
			user, err := apiClient.GetCurrentUser()
			if err != nil {
				_ = sessionStore.Clear()
				st.Auth = state.Auth{}
				apiClient.Token = ""
				st.Route = state.RouteWelcome
				w.Invalidate()
				return
			}

			st.Auth.Email = user.Email
			_ = sessionStore.Save(storage.Session{Token: st.Auth.Token, Email: st.Auth.Email})
			w.Invalidate()
		}()
	}

	firstFrame := true
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			if firstFrame {
				log.Printf("PassGO drawing first frame")
				firstFrame = false
			}
			gtx := app.NewContext(&ops, e)
			invalidate := false
			invalidateFunc := func() {
				w.Invalidate()
			}

			ui.ApplyTheme(th, st.DarkMode)

			paint.Fill(gtx.Ops, th.Palette.Bg)

			// Global navigation.
			if shell.VaultBtn.Clicked(gtx) {
				st.Nav = state.NavVault
				st.Route = state.RouteVaultList
			}
			if shell.SettingsBtn.Clicked(gtx) {
				st.Nav = state.NavSettings
				st.Route = state.RouteSettings
			}
			if shell.LogoutBtn.Clicked(gtx) {
				st.Auth = state.Auth{}
				st.Vaults = []state.Vault{}
				st.VaultsLoaded = false
				st.VaultsLoading = false
				st.VaultsLoadError = ""
				st.VaultsLoadDone = nil
				apiClient.Token = ""
				_ = sessionStore.Clear()
				st.Route = state.RouteWelcome
			}

			// Layout.
			switch st.Route {
			case state.RouteLogin:
				handleLoginPage(gtx, th, st, loginPage, apiClient, sessionStore, invalidateFunc)
			case state.RouteRegister:
				handleRegisterPage(gtx, th, st, registerPage, apiClient, invalidateFunc)
			case state.RouteForgotPassword:
				handleForgotPasswordPage(gtx, th, st, forgotPasswordPage, apiClient, invalidateFunc)
			case state.RouteVaultList:
				invalidate = handleVaultListPage(gtx, th, st, shell, vaultListPage, vaultAddPage, vaultDetailPage, apiClient, invalidateFunc)
			case state.RouteVaultAdd:
				invalidate = handleVaultAddPage(gtx, th, st, shell, vaultAddPage, apiClient, invalidateFunc)
			case state.RouteVaultDetail:
				invalidate = handleVaultDetailPage(gtx, th, st, shell, vaultDetailPage, apiClient, invalidateFunc)
			case state.RouteSettings:
				invalidate = handleSettingsPage(gtx, th, st, shell, settingsPage, apiClient, apiBaseURL, invalidateFunc)
			case state.RouteWelcome:
				fallthrough
			default:
				handleWelcomePage(gtx, th, st, welcomePage)
			}

			if invalidate {
				gtx.Execute(op.InvalidateCmd{})
			}
			e.Frame(gtx.Ops)
		}
	}
}
