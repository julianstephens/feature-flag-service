package commands

import (
	"context"
	"strings"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/feature-flag-service/internal/utils"
	"github.com/julianstephens/go-utils/cliutil"
	"github.com/julianstephens/go-utils/validator"
)

type UserCommand struct {
	Create struct {
		Email string `arg:"" help:"Email of the user to create."`
		Name  string `arg:"" help:"Name of the user to create."`
	} `cmd:"" help:"Create a new user."`
	List struct{} `cmd:"" help:"List all users."`
	Get  struct {
		ID    string `xor:"id,email" help:"ID of the user to get."`
		Email string `xor:"id,email" help:"Email of the user to get."`
	} `cmd:"" help:"Get a user."`
	Delete struct {
		ID    string `arg:"" help:"ID of the user to delete."`
		Email string `xor:"id,email" help:"Email of the user to delete."`
	} `cmd:"" help:"Delete a user."`
	Update struct {
		ID    string `arg:"" help:"ID of the user to update."`
		Email string `help:"New email of the user."`
		Name  string `help:"New name of the user."`
	} `cmd:"" help:"Update a user."`
}

func (c *UserCommand) ListUsers(ctx context.Context, conf *config.Config, client ffpb.RbacUserServiceClient) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DEFAULT_TIMEOUT)
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

func (c *UserCommand) GetUser(ctx context.Context, conf *config.Config, client ffpb.RbacUserServiceClient) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DEFAULT_TIMEOUT)
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

func (c *UserCommand) CreateUser(ctx context.Context, conf *config.Config, client ffpb.RbacUserServiceClient) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DEFAULT_TIMEOUT)
	defer cancel()

	password, err := utils.GenerateTempPassword()
	if err != nil {
		cliutil.PrintError("Failed to generate temporary password")
		return err
	}

	user, err := client.CreateUser(ctx, &ffpb.CreateUserRequest{
		Email:    c.Create.Email,
		Name:     c.Create.Name,
		Password: password,
	})
	if err != nil {
		cliutil.PrintError("Failed to create user")
		return err
	}

	utils.PrintUser(user.Id, user.Email, user.Name, user.CreatedAt, user.UpdatedAt, []string{})
	cliutil.PrintInfo("Temporary password: " + password + "\nPlease change this password after logging in.")

	return nil
}

func (c *UserCommand) UpdateUser(ctx context.Context, conf *config.Config, client ffpb.RbacUserServiceClient) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DEFAULT_TIMEOUT)
	defer cancel()

	if err := validator.ValidateUUID(c.Update.ID); err != nil {
		cliutil.PrintError("ID must be provided to update a user")
		return nil
	}

	user, err := client.GetUser(ctx, &ffpb.GetUserRequest{Id: c.Update.ID})
	if err != nil {
		cliutil.PrintError("Failed to fetch existing user")
		return err
	}

	req := &ffpb.UpdateUserRequest{
		Email: user.Email,
		Name:  user.Name,
	}

	if c.Update.Email != "" {
		if err := validator.ValidateEmail(c.Update.Email); err != nil {
			cliutil.PrintError("Invalid email format")
			return nil
		}
		req.Email = c.Update.Email
	}
	if c.Update.Name != "" {
		req.Name = c.Update.Name
	}

	_, err = client.UpdateUser(ctx, req)
	if err != nil {
		cliutil.PrintError("Failed to update user")
		return err
	}

	cliutil.PrintInfo("User updated successfully")
	return nil
}

func (c *UserCommand) DeleteUser(ctx context.Context, conf *config.Config, client ffpb.RbacUserServiceClient) error {
	ctx, cancel := context.WithTimeout(ctx, utils.DEFAULT_TIMEOUT)
	defer cancel()

	if c.Delete.ID != "" && c.Delete.Email != "" {
		cliutil.PrintError("Either ID or Email must be provided")
		return nil
	}

	var userID string
	if c.Delete.ID != "" {
		if err := validator.ValidateUUID(c.Delete.ID); err != nil {
			cliutil.PrintError("ID must be provided to delete a user")
			return nil
		}
		userID = c.Delete.ID
	} else if c.Delete.Email != "" {
		if err := validator.ValidateEmail(c.Delete.Email); err != nil {
			cliutil.PrintError("Invalid email format")
			return nil
		}
		user, err := client.GetUserByEmail(ctx, &ffpb.GetUserByEmailRequest{Email: c.Delete.Email})
		if err != nil {
			cliutil.PrintError("Failed to fetch user by email")
			return err
		}
		userID = user.Id
	}

	_, err := client.DeleteUser(ctx, &ffpb.DeleteUserRequest{Id: userID})
	if err != nil {
		cliutil.PrintError("Failed to delete user")
		return err
	}

	cliutil.PrintInfo("User deleted successfully")
	return nil
}
