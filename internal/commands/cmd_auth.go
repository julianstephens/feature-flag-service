package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/log"
	"google.golang.org/grpc"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/cache"
	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/utils"
	"github.com/julianstephens/go-utils/helpers"
	"github.com/julianstephens/go-utils/httputil/auth"
	"github.com/julianstephens/go-utils/jsonutil"
	"github.com/julianstephens/go-utils/security"
)

const (
	authFileName string = "auth.json"
	keyFileName  string = "key.bin"
)

type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type SecureAuthData struct {
	Credentials []byte `json:"credentials"`
}

type AuthData struct {
	Credentials Credentials
}

type AuthCommand struct {
	Login struct {
		Email string `arg:"" help:"Email of the user to log in as."`
	} `cmd:"" help:"Login to the featurectl CLI"`
	Status struct{} `cmd:"" help:"Check login status."`
}

func (c *AuthCommand) RunLogin(conf *config.Config, conn *grpc.ClientConn) error {
	validator := utils.Validator{}
	client := ffpb.NewAuthServiceClient(conn)
	req := &ffpb.LoginRequest{
		Email: c.Login.Email,
	}

	var password string
	huh.NewInput().EchoMode(huh.EchoModePassword).Validate(validator.ValidatePassword).Title("Password: ").Value(&password).Run()
	req.Password = password

	ctx, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
	defer cancel()

	err := spinner.New().Context(ctx).Title("Logging in...").ActionWithErr(func(ctx context.Context) error {
		return login(ctx, req, client)
	}).Run()
	if err != nil {
		return err
	}

	log.Info("Login successful")

	return nil
}

func (c *AuthCommand) RunStatus(mgr *auth.JWTManager) error {
	_, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
	defer cancel()

	cachePath, err := cache.Dir()
	if err != nil {
		log.Error("Error getting cache directory", "error", err)
		log.Info("Please set the XDG_CACHE_HOME environment variable to a writable directory")
		return err
	}

	if !helpers.Exists(cachePath+"/"+keyFileName) || !helpers.Exists(cachePath+"/"+authFileName) {
		log.Info("Not logged in. Use 'featurectl auth login' to log in.")
		return nil
	}

	key, err := loadKey()
	if err != nil {
		log.Warn("Malformed login cache, please log in again")
		cache.Remove(keyFileName)
		cache.Remove(authFileName)
		return err
	}

	authData, err := loadAuth(key)
	if err != nil || authData == nil {
		log.Warn("Malformed login cache, please log in again")
		cache.Remove(keyFileName)
		cache.Remove(authFileName)
		return err
	}
	expiresAt := authData.Credentials.ExpiresAt
	now := time.Now().Unix()
	if expiresAt > 0 && now < (expiresAt-60) {
		claims, err := mgr.ValidateToken(authData.Credentials.AccessToken)
		if err != nil {
			log.Warn("Invalid access token, please log in again")
			cache.Remove(keyFileName)
			cache.Remove(authFileName)
			return err
		}
		log.Info("Logged in", "user", claims.Email, "id", claims.Subject, "expires_at", time.Unix(expiresAt, 0).Format(time.RFC1123))
	} else {
		log.Info("Login expired, please log in again")
		cache.Remove(keyFileName)
		cache.Remove(authFileName)
	}

	return nil
}

func login(ctx context.Context, req *ffpb.LoginRequest, client ffpb.AuthServiceClient) (err error) {
	var creds Credentials

	// 1. Check for existing cache and load key
	key, err := loadKey()
	if err != nil {
		err = fmt.Errorf("error loading key: %w", err)
		return
	}

	// 2. If no cache, generate & save new key
	if len(key) == 0 {
		key, err = security.GenerateRandomKey(32)
		if err != nil {
			err = fmt.Errorf("error generating key: %w", err)
			return
		}
		if err = cache.WriteBytes(keyFileName, key); err != nil {
			err = fmt.Errorf("error writing key: %w", err)
			return
		}
	}

	// 3. Use key to decrypt existing cache if exists
	var authCache *AuthData
	authCache, err = loadAuth(key)
	if err != nil {
		return
	}

	var resp *ffpb.LoginResponse
	if authCache == nil {
		authCache = &AuthData{}
		// No existing cache, proceed to login
		resp, err = client.Login(ctx, req)
		if err != nil {
			err = fmt.Errorf("error logging in: %w", err)
			return
		}
		authCache.Credentials = Credentials{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			ExpiresAt:    time.Now().Unix() + resp.ExpiresIn,
		}
	} else {
		// Existing cache found, check if still valid
		expiresAt := authCache.Credentials.ExpiresAt
		if expiresAt > 0 && time.Now().Unix() < (expiresAt-60) {
			log.Info("Already logged in")
			return
		}
		// If expired, refresh tokens
		resp, err = client.Refresh(ctx, &ffpb.RefreshRequest{
			RefreshToken: authCache.Credentials.RefreshToken,
		})
		if err != nil {
			err = fmt.Errorf("error refreshing token: %w", err)
			return
		}
		authCache.Credentials = Credentials{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			ExpiresAt:    time.Now().Unix() + resp.ExpiresIn,
		}
	}
	creds = authCache.Credentials

	// 4. Encrypt and save cache
	if err = secureSave(key, creds); err != nil {
		return
	}
	return
}

func loadKey() (key []byte, err error) {
	key, err = cache.ReadBytes(keyFileName)
	if err != nil {
		return
	}
	return
}

func loadAuth(key []byte) (authData *AuthData, err error) {
	var secureCache SecureAuthData
	var parsedCreds Credentials
	var encrypted []byte

	authData = &AuthData{}

	err = cache.ReadJSON(authFileName, &secureCache)
	if err != nil {
		if !os.IsNotExist(err) {
			err = fmt.Errorf("error reading auth cache: %w", err)
			return
		}
		return nil, nil
	}

	encrypted = secureCache.Credentials
	if len(encrypted) > 0 {
		parsedCreds, err = decryptCredentials(key, encrypted)
		if err != nil {
			err = fmt.Errorf("error decrypting credentials: %w", err)
			return
		}
		authData.Credentials = parsedCreds
		return authData, nil
	}

	return nil, nil
}

func decryptCredentials(key []byte, encrypted []byte) (creds Credentials, err error) {
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

func secureSave(key []byte, creds Credentials) (err error) {
	secureCache := SecureAuthData{}
	encryptedResponse, err := security.Encrypt(key, helpers.MustMarshalJson(creds))
	if err != nil {
		return
	}
	secureCache.Credentials = encryptedResponse

	if err = cache.WriteJSON(authFileName, secureCache); err != nil {
		err = fmt.Errorf("error writing auth cache: %w", err)
		return
	}
	return
}
