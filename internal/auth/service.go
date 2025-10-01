package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/rbac/users"
	authutils "github.com/julianstephens/go-utils/httputil/auth"
	"github.com/julianstephens/go-utils/security"
	"github.com/julianstephens/go-utils/validator"
)

type Service interface {
	Login(ctx context.Context, email, password string) (*TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	Activate(ctx context.Context, email, password, newPassword string) (*TokenResponse, error)
}

type AuthClient struct {
	issuer  string
	userSvc *users.RbacUserService
	Manager *authutils.JWTManager
}

func NewAuthClient(conf *config.Config, userSvc *users.RbacUserService) (Service, error) {
	manager, err := authutils.NewJWTManager(conf.JWTSecret, time.Duration(conf.JWTExpiry)*time.Second, conf.JWTIssuer)
	if err != nil {
		return nil, err
	}
	return &AuthClient{
		issuer:  conf.JWTIssuer,
		userSvc: userSvc,
		Manager: manager,
	}, nil
}

func (a *AuthClient) Login(ctx context.Context, email, password string) (*TokenResponse, error) {
	if email == "" {
		return nil, fmt.Errorf("email must be provided")
	}
	if password == "" {
		return nil, fmt.Errorf("password must be provided")
	}

	user, err := a.userSvc.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	resp, err := a.Manager.GenerateTokenPairWithUserInfo(
		user.ID,
		user.Email,
		user.Email,
		user.Roles,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return &TokenResponse{
		Token:        resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    int(resp.ExpiresIn),
	}, nil
}

func (a *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token must be provided")
	}

	resp, err := a.Manager.ExchangeRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &TokenResponse{
		Token:        resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    int(resp.ExpiresIn),
	}, nil
}

func (a *AuthClient) Activate(ctx context.Context, email, password, newPassword string) (*TokenResponse, error) {
	if err := validator.ValidateNonEmpty(password); err != nil {
		return nil, fmt.Errorf("temporary password from account creation must be provided: %w", err)
	}
	if err := validator.ValidateEmail(email); err != nil {
		return nil, fmt.Errorf("invalid email format: %w", err)
	}
	if err := validator.ValidatePassword(newPassword); err != nil {
		return nil, fmt.Errorf("new password does not meet complexity requirements: %w", err)
	}

	user, err := a.userSvc.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if ok := security.VerifyPassword(password, user.Password); !ok {
		return nil, fmt.Errorf("invalid temporary password")
	}

	if err := a.userSvc.ActivateUser(ctx, user.ID, newPassword); err != nil {
		return nil, fmt.Errorf("failed to update password: %w", err)
	}

	if len(user.Roles) == 0 {
		user.Roles = []string{"user"}
	}

	tokenResp, err := a.Manager.GenerateTokenPairWithUserInfo(user.ID, user.Email, user.Email, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &TokenResponse{
		Token:        tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    int(tokenResp.ExpiresIn),
	}, nil
}
