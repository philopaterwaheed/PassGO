package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/philopaterwaheed/passGO/internal/backend/config"
)

var (
	ErrSupabaseNotConfigured = errors.New("supabase not configured")
	ErrSignupFailed          = errors.New("signup failed")
	ErrLoginFailed           = errors.New("login failed")
	ErrEmailNotVerified      = errors.New("email not verified")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrUserAlreadyExists     = errors.New("user already exists")
)

//handles communication with Supabase Auth
type SupabaseClient struct {
	url    string
	apiKey string
	client *http.Client
}

// SupabaseUser represents a user from Supabase Auth
type SupabaseUser struct {
	ID               string    `json:"id"`
	Email            string    `json:"email"`
	EmailConfirmedAt string    `json:"email_confirmed_at,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// SupabaseAuthResponse represents the response from Supabase Auth endpoints
type SupabaseAuthResponse struct {
	AccessToken  string       `json:"access_token,omitempty"`
	TokenType    string       `json:"token_type,omitempty"`
	ExpiresIn    int          `json:"expires_in,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	User         SupabaseUser `json:"user"`
}

// SupabaseErrorResponse represents an error response from Supabase
type SupabaseErrorResponse struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	Message          string `json:"msg,omitempty"`
	Code             int    `json:"code,omitempty"`
}

// NewSupabaseClient creates a new Supabase client
func NewSupabaseClient() (*SupabaseClient, error) {
	if config.SupabaseURL == "" || config.SupabaseAPIKey == "" {
		return nil, ErrSupabaseNotConfigured
	}

	return &SupabaseClient{
		url:    config.SupabaseURL,
		apiKey: config.SupabaseAPIKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

//registers a new user with Supabase Auth
func (s *SupabaseClient) SignUp(email, password string) (*SupabaseAuthResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.url+"/auth/v1/signup", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errResp SupabaseErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			if errResp.Message != "" {
				if errResp.Message == "User already registered" {
					return nil, ErrUserAlreadyExists
				}
				return nil, fmt.Errorf("%s", errResp.Message)
			}
			if errResp.ErrorDescription != "" {
				return nil, fmt.Errorf("%s", errResp.ErrorDescription)
			}
		}
		return nil, ErrSignupFailed
	}

	var authResp SupabaseAuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

// SignIn authenticates a user with email and password
func (s *SupabaseClient) SignIn(email, password string) (*SupabaseAuthResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.url+"/auth/v1/token?grant_type=password", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errResp SupabaseErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			if errResp.ErrorDescription == "Invalid login credentials" {
				return nil, ErrInvalidCredentials
			}
			if errResp.ErrorDescription == "Email not confirmed" {
				return nil, ErrEmailNotVerified
			}
			if errResp.ErrorDescription != "" {
				return nil, fmt.Errorf("%s", errResp.ErrorDescription)
			}
		}
		return nil, ErrLoginFailed
	}

	var authResp SupabaseAuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, err
	}

	// Check if email is verified
	if authResp.User.EmailConfirmedAt == "" {
		return nil, ErrEmailNotVerified
	}

	return &authResp, nil
}

// GetUser retrieves user information using an access token
func (s *SupabaseClient) GetUser(accessToken string) (*SupabaseUser, error) {
	req, err := http.NewRequest("GET", s.url+"/auth/v1/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New("failed to get user")
	}

	var user SupabaseUser
	if err := json.Unmarshal(respBody, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// ResendVerificationEmail resends the email verification link
func (s *SupabaseClient) ResendVerificationEmail(email string) error {
	payload := map[string]string{
		"email": email,
		"type":  "signup",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.url+"/auth/v1/resend", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("failed to resend verification email")
	}

	return nil
}

// VerifyOTP verifies an email OTP token
func (s *SupabaseClient) VerifyOTP(email, token, tokenType string) (*SupabaseAuthResponse, error) {
	payload := map[string]string{
		"email": email,
		"token": token,
		"type":  tokenType,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.url+"/auth/v1/verify", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New("verification failed")
	}

	var authResp SupabaseAuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

// ResetPassword sends a password reset email
func (s *SupabaseClient) ResetPassword(email string) error {
	payload := map[string]string{
		"email": email,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.url+"/auth/v1/recover", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("failed to send password reset email")
	}

	return nil
}

// UpdatePassword updates the user's password (requires valid access token)
func (s *SupabaseClient) UpdatePassword(accessToken, newPassword string) error {
	payload := map[string]string{
		"password": newPassword,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", s.url+"/auth/v1/user", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.apiKey)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("failed to update password")
	}

	return nil
}
