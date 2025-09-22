package commands

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"google.golang.org/grpc"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/utils"
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

func (c *FlagCommand) ListFlags(conf *config.Config, conn *grpc.ClientConn) error {
	client := ffpb.NewFlagServiceClient(conn)
	req := &ffpb.ListFlagsRequest{}

	res, err := client.ListFlags(context.Background(), req)
	if err != nil {
		log.Error("Failed to list flags")
		return err
	}

	if len(res.Flags) == 0 {
		log.Info("No flags found")
		return nil
	}

	var rows [][]string
	for _, flag := range res.Flags {
		rows = append(rows, []string{flag.Id, flag.Name, flag.Description, fmt.Sprintf("%v", flag.Enabled), flag.CreatedAt, flag.UpdatedAt})
	}

	utils.PrintTable([]string{"ID", "Name", "Description", "Enabled", "Created At", "Updated At"}, rows)

	return nil
}

func (c *FlagCommand) GetFlag(conf *config.Config, conn *grpc.ClientConn) error {
	client := ffpb.NewFlagServiceClient(conn)
	req := &ffpb.GetFlagRequest{
		Id: c.Get.ID,
	}

	flag, err := client.GetFlag(context.Background(), req)
	if err != nil {
		log.Error("Failed to get flag")
		return err
	}

	pprintFlag(flag)

	return nil
}

func (c *FlagCommand) CreateFlag(conf *config.Config, conn *grpc.ClientConn) error {
	client := ffpb.NewFlagServiceClient(conn)
	req := &ffpb.CreateFlagRequest{
		Name:        c.Create.Name,
		Description: c.Create.Description,
		Enabled:     c.Create.Enabled,
	}

	flag, err := client.CreateFlag(context.Background(), req)
	if err != nil {
		log.Error("Failed to create flag")
		return err
	}

	pprintFlag(flag)

	return nil
}

func (c *FlagCommand) UpdateFlag(conf *config.Config, conn *grpc.ClientConn) error {
	client := ffpb.NewFlagServiceClient(conn)
	req := &ffpb.UpdateFlagRequest{
		Id:      c.Update.ID,
		Enabled: c.Update.Enabled,
	}

	flag, err := client.GetFlag(context.Background(), &ffpb.GetFlagRequest{Id: c.Update.ID})
	if err != nil {
		log.Error("Failed to get existing flag")
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
		log.Error("Failed to update flag")
		return err
	}

	pprintFlag(flag)
	return nil
}

func (c *FlagCommand) DeleteFlag(conf *config.Config, conn *grpc.ClientConn) error {
	client := ffpb.NewFlagServiceClient(conn)
	req := &ffpb.DeleteFlagRequest{
		Id: c.Delete.ID,
	}

	_, err := client.DeleteFlag(context.Background(), req)
	if err != nil {
		log.Error("Failed to delete flag")
		return err
	}

	log.Info("Flag deleted successfully")
	return nil
}

func pprintFlag(flag *ffpb.Flag) {
	fmt.Printf("ID: %s\n", flag.Id)
	fmt.Printf("Name: %s\n", flag.Name)
	fmt.Printf("Description: %s\n", flag.Description)
	fmt.Printf("Enabled: %v\n", flag.Enabled)
	fmt.Printf("Created At: %s\n", flag.CreatedAt)
	fmt.Printf("Updated At: %s\n", flag.UpdatedAt)
}
