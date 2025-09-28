package users

import (
	"context"
	"time"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/rbac"
)

type RbacUserGRPCServer struct {
	ffpb.UnimplementedRbacUserServiceServer
	Service Service
}

func (s *RbacUserGRPCServer) ListUsers(ctx context.Context, req *ffpb.ListUsersRequest) (*ffpb.ListUsersResponse, error) {
	users, err := s.Service.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	var protoUsers []*ffpb.RbacUser
	for _, u := range users {
		protoUsers = append(protoUsers, ToProto(u))
	}
	return &ffpb.ListUsersResponse{Users: protoUsers}, nil
}

func (s *RbacUserGRPCServer) GetUser(ctx context.Context, req *ffpb.GetUserRequest) (*ffpb.RbacUser, error) {
	user, err := s.Service.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return ToProto(user), nil
}

func (s *RbacUserGRPCServer) GetUserByEmail(ctx context.Context, req *ffpb.GetUserByEmailRequest) (*ffpb.RbacUser, error) {
	user, err := s.Service.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return ToProto(user), nil
}

func (s *RbacUserGRPCServer) CreateUser(ctx context.Context, req *ffpb.CreateUserRequest) (*ffpb.RbacUser, error) {
	user, err := s.Service.CreateUser(ctx, req.Email, req.Name, req.Password, []string{})
	if err != nil {
		return nil, err
	}
	return ToProto(user), nil
}

func ToProto(u *rbac.RbacUserDto) *ffpb.RbacUser {
	return &ffpb.RbacUser{
		Id:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		Roles:     u.Roles,
	}
}
