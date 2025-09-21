package auth

import (
	"time"

	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/go-utils/helpers"
	authutils "github.com/julianstephens/go-utils/httputil/auth"
)

type AuthClient struct {
	issuer string
	manager *authutils.JWTManager
}


func NewAuthClient(conf *config.Config) *AuthClient {
	return &AuthClient{
		issuer: conf.JWTIssuer,
		manager: authutils.NewJWTManager(conf.JWTSecret, time.Duration(conf.JWTExpiry) * time.Second, conf.JWTIssuer),
	}
}

func (a *AuthClient) Authenticate(token string) (bool, error) {
	claims, err := a.manager.ValidateToken(token)
	if err != nil {
		return false, err
	}
	return a.validateClaims(claims), nil
}

func (a *AuthClient) GenerateToken(userID string, roles []string, customClaims *map[string]any) (string, error) {
	return a.manager.GenerateTokenWithClaims(userID, roles, helpers.Default(*customClaims, make(map[string]any, 0)))
}

func (a *AuthClient) validateClaims(claims *authutils.Claims) bool {
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