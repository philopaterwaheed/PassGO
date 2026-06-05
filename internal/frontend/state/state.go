package state

import (
	"fmt"
	"time"
)

// Route represents the current screen.
type Route int

const (
	RouteWelcome Route = iota
	RouteLogin
	RouteRegister
	RouteForgotPassword
	RouteVaultList
	RouteVaultAdd
	RouteVaultDetail
	RouteSettings
)

// NavItem is a top-level destination shown in the app chrome.
type NavItem int

const (
	NavVault NavItem = iota
	NavSettings
)

type Auth struct {
	Token          string
	Email          string
	MasterPassword string
}

// Vault represents a single credential entry stored
type Vault struct {
	ID       string
	Title    string
	Username string
	Password string
	URL      string
	Notes    string
}

type AppState struct {
	Route Route
	Nav   NavItem
	Auth  Auth

	DarkMode bool

	Vaults          []Vault
	SelectedVaultID string
	VaultsLoaded    bool
}

func (s *AppState) IsAuthed() bool {
	return s.Auth.Token != ""
}

func (s *AppState) AddVault(v Vault) {
	if v.ID == "" {
		v.ID = NewVaultID()
	}
	s.Vaults = append([]Vault{v}, s.Vaults...)
}

func (s *AppState) VaultByID(id string) (Vault, bool) {
	for _, v := range s.Vaults {
		if v.ID == id {
			return v, true
		}
	}
	return Vault{}, false
}

func NewVaultID() string {
	return fmt.Sprintf("v_%d", time.Now().UnixNano())
}
