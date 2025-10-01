package auth

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"google.golang.org/grpc/metadata"

	"github.com/julianstephens/feature-flag-service/internal/cache"
	"github.com/julianstephens/feature-flag-service/internal/types"
	"github.com/julianstephens/feature-flag-service/internal/utils"
	"github.com/julianstephens/go-utils/helpers"
	"github.com/julianstephens/go-utils/httputil/auth"
	"github.com/julianstephens/go-utils/jsonutil"
	"github.com/julianstephens/go-utils/security"
)

// TokenManager handles loading and managing authentication tokens
type TokenManager struct {
	jwtManager *auth.JWTManager
}

// NewTokenManager creates a new token manager
func NewTokenManager(jwtManager *auth.JWTManager) *TokenManager {
	return &TokenManager{
		jwtManager: jwtManager,
	}
}

// LoadToken attempts to load a valid token from cache
func (tm *TokenManager) LoadToken() (string, error) {
	cachePath, err := cache.Dir()
	if err != nil {
		return "", err
	}

	// Check if auth files exist
	if !helpers.Exists(cachePath+"/"+utils.DEFAULT_KEY_FILE) || !helpers.Exists(cachePath+"/"+utils.DEFAULT_AUTH_CACHE_FILE) {
		return "", ErrNotLoggedIn
	}

	// Load encryption key
	key, err := LoadKey()
	if err != nil {
		log.Warn("Malformed login cache, please log in again")
		cache.Remove(utils.DEFAULT_KEY_FILE)
		cache.Remove(utils.DEFAULT_AUTH_CACHE_FILE)
		return "", ErrNotLoggedIn
	}

	// Load encrypted auth data
	var authData *types.AuthData
	authData, err = LoadAuth(key)
	if err != nil {
		log.Warn("Malformed login cache, please log in again")
		cache.Remove(utils.DEFAULT_KEY_FILE)
		cache.Remove(utils.DEFAULT_AUTH_CACHE_FILE)
		return "", ErrNotLoggedIn
	}

	// Validate token
	_, err = tm.jwtManager.ValidateToken(authData.Credentials.AccessToken)
	if err != nil {
		cache.Remove(utils.DEFAULT_KEY_FILE)
		cache.Remove(utils.DEFAULT_AUTH_CACHE_FILE)
		return "", ErrTokenExpired
	}

	return authData.Credentials.AccessToken, nil
}

// CreateAuthenticatedContext creates a gRPC context with authentication metadata
func (tm *TokenManager) CreateAuthenticatedContext(ctx context.Context) (context.Context, error) {
	token, err := tm.LoadToken()
	if err != nil {
		return ctx, err
	}

	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	return metadata.NewOutgoingContext(ctx, md), nil
}

// Custom errors
var (
	ErrNotLoggedIn  = &AuthError{"not logged in - please run 'featurectl auth login'"}
	ErrTokenExpired = &AuthError{"token expired - please run 'featurectl auth login'"}
)

type AuthError struct {
	message string
}

func (e *AuthError) Error() string {
	return e.message
}

func LoadKey() (key []byte, err error) {
	key, err = cache.ReadBytes(utils.DEFAULT_KEY_FILE)
	if err != nil {
		return
	}
	return
}

func LoadAuth(key []byte) (authData *types.AuthData, err error) {
	var secureCache types.SecureAuthData
	var parsedCreds types.Credentials
	var encrypted []byte

	authData = &types.AuthData{}

	err = cache.ReadJSON(utils.DEFAULT_AUTH_CACHE_FILE, &secureCache)
	if err != nil {
		if !os.IsNotExist(err) {
			err = fmt.Errorf("error reading auth cache: %w", err)
			return
		}
		return nil, nil
	}

	encrypted = secureCache.Credentials
	if len(encrypted) > 0 {
		parsedCreds, err = DecryptCredentials(key, encrypted)
		if err != nil {
			err = fmt.Errorf("error decrypting credentials: %w", err)
			return
		}
		authData.Credentials = parsedCreds
		return authData, nil
	}

	return nil, nil
}

func DecryptCredentials(key []byte, encrypted []byte) (creds types.Credentials, err error) {
	var decrypted []byte
	decrypted, err = security.Decrypt(key, encrypted)
	if err != nil {
		return
	}

	if err = jsonutil.DecodeReader(bytes.NewReader(decrypted), &creds, &jsonutil.DecoderOptions{}); err != nil {
		return
	}

	return
}

func SecureSave(key []byte, creds types.Credentials) (err error) {
	secureCache := types.SecureAuthData{}
	encryptedResponse, err := security.Encrypt(key, helpers.MustMarshalJson(creds))
	if err != nil {
		return
	}
	secureCache.Credentials = encryptedResponse

	if err = cache.WriteJSON(utils.DEFAULT_AUTH_CACHE_FILE, secureCache); err != nil {
		err = fmt.Errorf("error writing auth cache: %w", err)
		return
	}
	return
}
