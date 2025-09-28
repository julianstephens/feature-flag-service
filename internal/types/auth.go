package types

// Credentials represents user authentication credentials
type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// AuthData represents the decrypted authentication data
type AuthData struct {
	Credentials Credentials `json:"credentials"`
}

// SecureAuthData represents the encrypted authentication data stored in cache
type SecureAuthData struct {
	Credentials []byte `json:"credentials"`
}
