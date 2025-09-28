package flag

import (
	"context"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
)

type FlagGRPCServer struct {
	ffpb.UnimplementedFlagServiceServer
	Service Service
}

func (s *FlagGRPCServer) ListFlags(ctx context.Context, req *ffpb.ListFlagsRequest) (*ffpb.ListFlagsResponse, error) {
	flags, err := s.Service.ListFlags(ctx)
	if err != nil {
		return nil, err
	}
	var protoFlags []*ffpb.Flag
	for _, f := range flags {
		protoFlags = append(protoFlags, f.ToProto())
	}
	return &ffpb.ListFlagsResponse{Flags: protoFlags}, nil
}

func (s *FlagGRPCServer) GetFlag(ctx context.Context, req *ffpb.GetFlagRequest) (*ffpb.Flag, error) {
	flag, err := s.Service.GetFlag(ctx, req.Id)

	if err != nil {
		return nil, err
	}
	return flag.ToProto(), nil
}

func (s *FlagGRPCServer) CreateFlag(ctx context.Context, req *ffpb.CreateFlagRequest) (*ffpb.Flag, error) {
	flag, err := s.Service.CreateFlag(ctx, req.Name, req.Description, req.Enabled)
	if err != nil {
		return nil, err
	}
	return flag.ToProto(), nil
}

func (s *FlagGRPCServer) UpdateFlag(ctx context.Context, req *ffpb.UpdateFlagRequest) (*ffpb.Flag, error) {
	flag, err := s.Service.UpdateFlag(ctx, req.Id, req.Name, req.Description, req.Enabled)
	if err != nil {
		return nil, err
	}
	return flag.ToProto(), nil
}

func (s *FlagGRPCServer) DeleteFlag(ctx context.Context, req *ffpb.DeleteFlagRequest) (*ffpb.DeleteFlagResponse, error) {
	err := s.Service.DeleteFlag(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &ffpb.DeleteFlagResponse{}, nil
}
