package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/auth"
	"github.com/julianstephens/feature-flag-service/internal/cache"
	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/types"
	"github.com/julianstephens/feature-flag-service/internal/utils"
	"github.com/julianstephens/go-utils/cliutil"
	authutil "github.com/julianstephens/go-utils/httputil/auth"
	"github.com/julianstephens/go-utils/security"
)

const (
	keyFileName string = "key.bin"
)

var (
	ErrAlreadyLoggedIn = fmt.Errorf("already logged in")
)

type AuthCommand struct {
	Login struct {
		Email string `arg:"" help:"Email of the user to log in as."`
	} `cmd:"" help:"Login to the featurectl CLI"`
	Status struct{} `cmd:"" help:"Check login status."`
}

func (c *AuthCommand) RunLogin(conf *config.Config, conn *grpc.ClientConn) error {
	client := ffpb.NewAuthServiceClient(conn)
	req := &ffpb.LoginRequest{
		Email: c.Login.Email,
	}

	password := cliutil.PromptPassword("Password: ")
	req.Password = password

	ctx, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
	defer cancel()

	spinner := cliutil.NewSpinner("Logging in...")
	spinner.Start()
	err := login(ctx, req, client)
	time.Sleep(500 * time.Millisecond) // Ensure spinner shows for at least half a second
	spinner.Stop()
	if err != nil && !errors.Is(err, ErrAlreadyLoggedIn) {
		return err
	} else if err != nil && errors.Is(err, ErrAlreadyLoggedIn) {
		cliutil.PrintInfo("Already logged in")
	} else {
		cliutil.PrintSuccess("Login successful")
	}

	return nil
}

func (c *AuthCommand) RunStatus(ctx context.Context, mgr *authutil.JWTManager) error {
	_, cancel := context.WithTimeout(ctx, utils.DEFAULT_TIMEOUT)
	defer cancel()

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		cliutil.PrintInfo("Not logged in. Use 'featurectl auth login' to log in.")
		return nil
	}
	if len(md["authorization"]) == 0 {
		cliutil.PrintInfo("Not logged in. Use 'featurectl auth login' to log in.")
		return nil
	}

	token := strings.Split(md["authorization"][0], "Bearer ")[1]

	claims, err := mgr.ValidateToken(token)
	if err != nil {
		cliutil.PrintInfo("Not logged in. Use 'featurectl auth login' to log in.")
		return nil
	}

	cliutil.PrintSuccess("Logged in")
	utils.PrintUser(claims.Subject, claims.Email, "N/A", time.Unix(claims.IssuedAt.Unix(), 0).Format(time.RFC1123), time.Unix(claims.ExpiresAt.Unix(), 0).Format(time.RFC1123), claims.Roles)
	return nil
}

func login(ctx context.Context, req *ffpb.LoginRequest, client ffpb.AuthServiceClient) (err error) {
	var creds types.Credentials

	// 1. Check for existing cache and load key
	key, err := auth.LoadKey()
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
	var authCache *types.AuthData
	authCache, err = auth.LoadAuth(key)
	if err != nil {
		return
	}

	var resp *ffpb.LoginResponse
	if authCache == nil {
		authCache = &types.AuthData{}
		// No existing cache, proceed to login
		resp, err = client.Login(ctx, req)
		if err != nil {
			err = fmt.Errorf("error logging in: %w", err)
			return
		}
		authCache.Credentials = types.Credentials{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			ExpiresAt:    time.Now().Unix() + resp.ExpiresIn,
		}
	} else {
		// Existing cache found, check if still valid
		expiresAt := authCache.Credentials.ExpiresAt
		if expiresAt > 0 && time.Now().Unix() < (expiresAt-60) {
			return ErrAlreadyLoggedIn
		}
		// If expired, refresh tokens
		resp, err = client.Refresh(ctx, &ffpb.RefreshRequest{
			RefreshToken: authCache.Credentials.RefreshToken,
		})
		if err != nil {
			err = fmt.Errorf("error refreshing token: %w", err)
			return
		}
		authCache.Credentials = types.Credentials{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			ExpiresAt:    time.Now().Unix() + resp.ExpiresIn,
		}
	}
	creds = authCache.Credentials

	// 4. Encrypt and save cache
	if err = auth.SecureSave(key, creds); err != nil {
		return
	}
	return
}
