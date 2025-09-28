package commands

import (
	"context"
	"strings"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/utils"
	"github.com/julianstephens/go-utils/cliutil"
)

type UserCommand struct {
	Create struct {
		Email string `arg:"" help:"Email of the user to create."`
	} `cmd:"" help:"Create a new user."`
	List struct{} `cmd:"" help:"List all users."`
	Get  struct {
		ID    string `xor:"id,email" help:"ID of the user to get."`
		Email string `xor:"id,email" help:"Email of the user to get."`
	} `cmd:"" help:"Get a user."`
	Delete struct {
		ID string `arg:"" help:"ID of the user to delete."`
	} `cmd:"" help:"Delete a user."`
	Update struct {
		ID    string   `arg:"" help:"ID of the user to update."`
		Email string   `help:"New email of the user."`
		Name  string   `help:"New name of the user."`
		Roles []string `help:"New roles for the user."`
	} `cmd:"" help:"Update a user."`
}

func (c *UserCommand) ListUsers(conf *config.Config, client ffpb.RbacUserServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
	defer cancel()

	users, err := client.ListUsers(ctx, &ffpb.ListUsersRequest{})
	if err != nil {
		cliutil.PrintError("Failed to list users")
		return err
	}

	var rows [][]string
	rows = append(rows, []string{"ID", "Email", "Name", "Created At", "Updated At", "Roles"})
	for _, user := range users.Users {
		rows = append(rows, []string{user.Id, user.Email, user.Name, user.CreatedAt, user.UpdatedAt, strings.Join(user.Roles, ",")})
	}
	cliutil.PrintTable(rows)

	return nil
}

func (c *UserCommand) GetUser(conf *config.Config, client ffpb.RbacUserServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), utils.DEFAULT_TIMEOUT)
	defer cancel()

	var user *ffpb.RbacUser
	var err error

	if c.Get.ID != "" {
		user, err = client.GetUser(ctx, &ffpb.GetUserRequest{Id: c.Get.ID})
	} else if c.Get.Email != "" {
		user, err = client.GetUserByEmail(ctx, &ffpb.GetUserByEmailRequest{Email: c.Get.Email})
	} else {
		cliutil.PrintError("Either ID or Email must be provided")
		return nil
	}

	if err != nil {
		cliutil.PrintError("Failed to get user")
		return err
	}

	utils.PrintUser(user.Id, user.Email, user.Name, user.CreatedAt, user.UpdatedAt, user.Roles)

	return nil
}

func (c *UserCommand) CreateUser(conf *config.Config, client ffpb.RbacUserServiceClient) error {
	return nil
}

func (c *UserCommand) UpdateUser(conf *config.Config, client ffpb.RbacUserServiceClient) error {
	return nil
}

func (c *UserCommand) DeleteUser(conf *config.Config, client ffpb.RbacUserServiceClient) error {
	return nil
}
