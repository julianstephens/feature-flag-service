package auth


type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn string `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}