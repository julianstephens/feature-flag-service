package commands

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/auth"
	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/utils"
	"github.com/julianstephens/go-utils/cliutil"
	authutils "github.com/julianstephens/go-utils/httputil/auth"
)

type FlagCommand struct {
	List struct{} `cmd:"" help:"List all feature flags."`
	Get  struct {
		ID string `arg:"" help:"ID of the feature flag to retrieve."`
	} `cmd:"" help:"Get details of a specific feature flag by ID."`
	Create struct {
		Name        string `help:"Name of the feature flag."`
		Description string `help:"Description of the feature flag."`
		Enabled     bool   `negatable:"disabled" help:"Initial state of the feature flag."`
	} `cmd:"" help:"Create a new feature flag."`
	Update struct {
		ID          string `arg:"" help:"ID of the feature flag to update."`
		Name        string `optional:"" help:"New name of the feature flag."`
		Description string `optional:"" help:"New description of the feature flag."`
		Enabled     bool   `negatable:"disabled" help:"New state of the feature flag."`
	} `cmd:"" help:"Update an existing feature flag by ID."`
	Delete struct {
		ID string `arg:"" help:"ID of the feature flag to delete."`
	} `cmd:"" help:"Delete a feature flag by ID."`
}

func (c *FlagCommand) ListFlags(conf *config.Config, conn *grpc.ClientConn, jwtManager *authutils.JWTManager) error {
	client := ffpb.NewFlagServiceClient(conn)
	req := &ffpb.ListFlagsRequest{}

	ctx, err := createAuthenticatedContext(jwtManager)
	if err != nil {
		cliutil.PrintError("Failed to create authenticated context")
		return err
	}

	cancelCtx, cancel := context.WithTimeout(ctx, utils.DEFAULT_TIMEOUT)
	defer cancel()

	res, err := client.ListFlags(cancelCtx, req)
	if err != nil {
		cliutil.PrintError("Failed to list flags")
		return err
	}

	if len(res.Flags) == 0 {
		cliutil.PrintInfo("No flags found")
		return nil
	}

	var rows [][]string
	rows = append(rows, []string{"ID", "Name", "Description", "Enabled", "Created At", "Updated At"})
	for _, flag := range res.Flags {
		rows = append(rows, []string{flag.Id, flag.Name, flag.Description, fmt.Sprintf("%v", flag.Enabled), flag.CreatedAt, flag.UpdatedAt})
	}

	cliutil.PrintTable(rows)

	return nil
}

func (c *FlagCommand) GetFlag(conf *config.Config, conn *grpc.ClientConn, jwtManager *authutils.JWTManager) error {
	client := ffpb.NewFlagServiceClient(conn)
	req := &ffpb.GetFlagRequest{
		Id: c.Get.ID,
	}

	flag, err := client.GetFlag(context.Background(), req)
	if err != nil {
		cliutil.PrintError("Failed to get flag")
		return err
	}

	pprintFlag(flag)

	return nil
}

func (c *FlagCommand) CreateFlag(conf *config.Config, conn *grpc.ClientConn, jwtManager *authutils.JWTManager) error {
	client := ffpb.NewFlagServiceClient(conn)
	req := &ffpb.CreateFlagRequest{
		Name:        c.Create.Name,
		Description: c.Create.Description,
		Enabled:     c.Create.Enabled,
	}

	flag, err := client.CreateFlag(context.Background(), req)
	if err != nil {
		cliutil.PrintError("Failed to create flag")
		return err
	}

	pprintFlag(flag)

	return nil
}

func (c *FlagCommand) UpdateFlag(conf *config.Config, conn *grpc.ClientConn, jwtManager *authutils.JWTManager) error {
	client := ffpb.NewFlagServiceClient(conn)
	req := &ffpb.UpdateFlagRequest{
		Id:      c.Update.ID,
		Enabled: c.Update.Enabled,
	}

	flag, err := client.GetFlag(context.Background(), &ffpb.GetFlagRequest{Id: c.Update.ID})
	if err != nil {
		cliutil.PrintError("Failed to get existing flag")
		return err
	}

	// Only update fields that were provided
	if c.Update.Name != flag.Name && c.Update.Name != "" {
		req.Name = c.Update.Name
	} else {
		req.Name = flag.Name
	}
	if c.Update.Description != flag.Description && c.Update.Description != "" {
		req.Description = c.Update.Description
	} else {
		req.Description = flag.Description
	}

	flag, err = client.UpdateFlag(context.Background(), req)
	if err != nil {
		cliutil.PrintError("Failed to update flag")
		return err
	}

	pprintFlag(flag)
	return nil
}

func (c *FlagCommand) DeleteFlag(conf *config.Config, conn *grpc.ClientConn, jwtManager *authutils.JWTManager) error {
	client := ffpb.NewFlagServiceClient(conn)
	req := &ffpb.DeleteFlagRequest{
		Id: c.Delete.ID,
	}

	_, err := client.DeleteFlag(context.Background(), req)
	if err != nil {
		cliutil.PrintError("Failed to delete flag")
		return err
	}

	cliutil.PrintInfo("Flag deleted successfully")
	return nil
}

func pprintFlag(flag *ffpb.Flag) {
	cliutil.PrintInfo(fmt.Sprintf("ID: %s\n", flag.Id))
	cliutil.PrintInfo(fmt.Sprintf("Name: %s\n", flag.Name))
	cliutil.PrintInfo(fmt.Sprintf("Description: %s\n", flag.Description))
	cliutil.PrintInfo(fmt.Sprintf("Enabled: %v\n", flag.Enabled))
	cliutil.PrintInfo(fmt.Sprintf("Created At: %s\n", flag.CreatedAt))
	cliutil.PrintInfo(fmt.Sprintf("Updated At: %s\n", flag.UpdatedAt))
}

// createAuthenticatedContext loads token using TokenManager and creates an authenticated gRPC context
func createAuthenticatedContext(jwtManager *authutils.JWTManager) (context.Context, error) {
	tokenManager := auth.NewTokenManager(jwtManager)
	ctx, err := tokenManager.CreateAuthenticatedContext(context.Background())
	if err != nil {
		if err.Error() == "invalid key size" || err.Error() == "key not found" {
			return nil, fmt.Errorf("authentication required - please run 'featurectl auth login <email>' first")
		}
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	return ctx, nil
}
