package models

// VerifyHashRequest represents the request to verify using an access token received in hash
type VerifyHashRequest struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
}
