package auth

import (
	"time"

	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/go-utils/helpers"
	authutils "github.com/julianstephens/go-utils/httputil/auth"
)

type AuthClient struct {
	issuer  string
	Manager *authutils.JWTManager
}

func NewAuthClient(conf *config.Config) (*AuthClient, error) {
	manager, err := authutils.NewJWTManager(conf.JWTSecret, time.Duration(conf.JWTExpiry)*time.Second, conf.JWTIssuer)
	if err != nil {
		return nil, err
	}
	return &AuthClient{
		issuer:  conf.JWTIssuer,
		Manager: manager,
	}, nil
}

// Authenticate validates the token and checks that the user ID in the claims matches the provided user ID.
func (a *AuthClient) Authenticate(token, userId string) (bool, error) {
	claims, err := a.Manager.ValidateToken(token)
	if err != nil {
		return false, err
	}
	return a.validateClaims(claims, userId), nil
}

// Issue generates a new JWT token and refresh token for the given user ID and roles.
func (a *AuthClient) Issue(userId string, roles []string, claims *map[string]any) (*TokenResponse, error) {
	token, err := a.GenerateToken(userId, roles, claims)
	if err != nil {
		return nil, err
	}

	exp, err := a.Manager.TokenExpiration(token)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		Token:        token,
		RefreshToken: token,
		ExpiresIn:    int(time.Until(exp).Seconds()),
	}, nil
}

// Refresh generates a new JWT token and refresh token using the provided refresh token.
func (a *AuthClient) Refresh(token string) (*TokenResponse, error) {
	token, err := a.Manager.RefreshToken(token)
	if err != nil {
		return nil, err
	}

	exp, err := a.Manager.TokenExpiration(token)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		Token:        token,
		RefreshToken: token,
		ExpiresIn:    int(time.Until(exp).Seconds()),
	}, nil
}

// GenerateToken generates a JWT token for the given user ID, roles, and custom claims.
func (a *AuthClient) GenerateToken(userID string, roles []string, customClaims *map[string]any) (string, error) {
	return a.Manager.GenerateTokenWithClaims(userID, roles, helpers.Default(*customClaims, make(map[string]any, 0)))
}

func (a *AuthClient) validateClaims(claims *authutils.Claims, userID string) bool {
	if userID != "" && claims.RegisteredClaims.Subject != userID {
		return false
	}
	if claims.RegisteredClaims.Issuer != a.issuer {
		return false
	}
	if claims.RegisteredClaims.ExpiresAt.Time.Before(time.Now()) {
		return false
	}
	if claims.RegisteredClaims.NotBefore.Time.After(time.Now()) {
		return false
	}
	if claims.RegisteredClaims.IssuedAt.Time.After(time.Now()) {
		return false
	}
	return true
}
