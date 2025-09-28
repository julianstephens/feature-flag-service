package auth

import (
	"context"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
)

type AuthGRPCServer struct {
	ffpb.UnimplementedAuthServiceServer
	Service Service
}

func (s *AuthGRPCServer) Login(ctx context.Context, req *ffpb.LoginRequest) (*ffpb.LoginResponse, error) {
	tokenResp, err := s.Service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &ffpb.LoginResponse{
		AccessToken:  tokenResp.Token,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    int64(tokenResp.ExpiresIn),
	}, nil
}

func (s *AuthGRPCServer) Refresh(ctx context.Context, req *ffpb.RefreshRequest) (*ffpb.LoginResponse, error) {
	tokenResp, err := s.Service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &ffpb.LoginResponse{
		AccessToken:  tokenResp.Token,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    int64(tokenResp.ExpiresIn),
	}, nil
}

func (s *AuthGRPCServer) Activate(ctx context.Context, req *ffpb.ActivateRequest) (*ffpb.LoginResponse, error) {
	tokenResp, err := s.Service.Activate(ctx, req.Email, req.Password, req.NewPassword)
	if err != nil {
		return nil, err
	}
	return &ffpb.LoginResponse{
		AccessToken:  tokenResp.Token,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    int64(tokenResp.ExpiresIn),
	}, nil
}
