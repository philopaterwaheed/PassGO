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
	ID          string
	Title       string
	Username    string
	Password    string
	URL         string
	Notes       string
	Decrypted   bool
	HasPassword bool
}

type AppState struct {
	Route Route
	Nav   NavItem
	Auth  Auth

	DarkMode bool

	Vaults          []Vault
	SelectedVaultID string
	VaultsLoaded    bool
	VaultsLoading   bool
	VaultsLoadError string
	VaultsLoadDone  chan VaultLoadResult

	VaultDetailLoading   bool
	VaultDetailLoadError string
	VaultDetailLoadDone  chan VaultDetailLoadResult
}

type VaultLoadResult struct {
	Vaults []Vault
	Err    error
}

type VaultDetailLoadResult struct {
	Vault Vault
	Err   error
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

func (s *AppState) UpdateVault(v Vault) {
	for i, existing := range s.Vaults {
		if existing.ID == v.ID {
			s.Vaults[i] = v
			return
		}
	}
	s.Vaults = append([]Vault{v}, s.Vaults...)
}

func (s *AppState) DeleteVault(id string) {
	for i, v := range s.Vaults {
		if v.ID == id {
			s.Vaults = append(s.Vaults[:i], s.Vaults[i+1:]...)
			if s.SelectedVaultID == id {
				s.SelectedVaultID = ""
			}
			return
		}
	}
}

func NewVaultID() string {
	return fmt.Sprintf("v_%d", time.Now().UnixNano())
}
