package grpcmiddleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/julianstephens/go-utils/httputil/auth"
)

// JWTAuthInterceptor validates JWT tokens from gRPC metadata
func JWTAuthInterceptor(jwtManager *auth.JWTManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		token := strings.TrimPrefix(authHeaders[0], "Bearer ")
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Add claims to context for handlers to use
		ctx = context.WithValue(ctx, "user_claims", claims)

		return handler(ctx, req)
	}
}

// RequireRolesInterceptor checks if user has required roles
func RequireRolesInterceptor(requiredRoles ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		claims, ok := ctx.Value("user_claims").(*auth.Claims)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "no user claims")
		}

		hasRole := false
		for _, userRole := range claims.Roles {
			for _, reqRole := range requiredRoles {
				if userRole == reqRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
		}

		return handler(ctx, req)
	}
}

// MethodBasedAuthInterceptor provides method-specific role requirements
func MethodBasedAuthInterceptor(jwtManager *auth.JWTManager, roleMap map[string][]string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Check if method requires auth
		requiredRoles, requiresAuth := roleMap[info.FullMethod]
		if !requiresAuth {
			return handler(ctx, req) // Public method
		}

		// Extract and validate JWT token
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		token := strings.TrimPrefix(authHeaders[0], "Bearer ")
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Check roles
		hasRole := false
		for _, userRole := range claims.Roles {
			for _, reqRole := range requiredRoles {
				if userRole == reqRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
		}

		// Add claims to context for handlers to use
		ctx = context.WithValue(ctx, "user_claims", claims)

		return handler(ctx, req)
	}
}
