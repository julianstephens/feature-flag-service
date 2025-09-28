package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ffpb "github.com/julianstephens/feature-flag-service/gen/go/grpc/v1/featureflag.v1"
	"github.com/julianstephens/feature-flag-service/internal/auth"
	"github.com/julianstephens/feature-flag-service/internal/commands"
	"github.com/julianstephens/feature-flag-service/internal/config"
	"github.com/julianstephens/go-utils/cliutil"
	authutils "github.com/julianstephens/go-utils/httputil/auth"
)

type Globals struct {
	Version kong.VersionFlag
}

type CLI struct {
	Globals

	Auth  commands.AuthCommand `cmd:"" help:"Manage authentication."`
	Flag  commands.FlagCommand `cmd:"" help:"Manage feature flags."`
	User  commands.UserCommand `cmd:"" help:"Manage users."`
	Audit struct {
	} `cmd:"" help:"Audit log operations."`
}

var err error

func main() {
	cli := CLI{
		Globals: Globals{},
	}

	conf := config.LoadConfig()
	var jwtManager *authutils.JWTManager
	jwtManager, err = authutils.NewJWTManager(conf.JWTSecret, time.Duration(conf.JWTExpiry), conf.JWTIssuer)
	if err != nil {
		log.Fatal("Failed to create JWT manager:", "error", err)
	}

	var conn *grpc.ClientConn
	conn, err = grpc.NewClient(":"+conf.GRPCPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to gRPC server:", "error", err)
	}
	defer conn.Close()

	kongCtx := kong.Parse(&cli,
		kong.Name("featurectl"),
		kong.Description("CLI for Distributed Feature Flag & Config System"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true, NoExpandSubcommands: true}),
		kong.Vars{"version": "1.0.0"},
	)

	ctx, err := createAuthenticatedContext(jwtManager)
	if err != nil {
		cliutil.PrintWarning(fmt.Sprintf("Failed to create authenticated context: %v", err))
	}

	cmd := strings.Split(kongCtx.Command(), " ")
	switch cmd[0] {
	case "auth":
		if len(cmd) < 2 {
			kongCtx.PrintUsage(false)
			return
		}
		subcmd := cmd[1]

		switch subcmd {
		case "login":
			err = cli.Auth.RunLogin(conf, conn)
		case "status":
			err = cli.Auth.RunStatus(ctx, jwtManager)
		default:
			panic(fmt.Sprintf("unknown auth command: %s", subcmd))
		}
	case "audit":
		// Implement audit log functionality here
	case "flag":
		if len(cmd) < 2 {
			kongCtx.PrintUsage(false)
			return
		}
		subcmd := cmd[1]
		switch subcmd {
		case "list":
			err = cli.Flag.ListFlags(conf, conn, jwtManager)
		case "get":
			err = cli.Flag.GetFlag(conf, conn, jwtManager)
		case "create":
			err = cli.Flag.CreateFlag(conf, conn, jwtManager)
		case "update":
			err = cli.Flag.UpdateFlag(conf, conn, jwtManager)
		case "delete":
			err = cli.Flag.DeleteFlag(conf, conn, jwtManager)
		default:
			panic(fmt.Sprintf("unknown flag command: %s", subcmd))
		}
	case "user":
		if len(cmd) < 2 {
			kongCtx.PrintUsage(false)
			return
		}
		client := ffpb.NewRbacUserServiceClient(conn)
		subcmd := cmd[1]
		switch subcmd {
		case "create":
			err = cli.User.CreateUser(conf, client)
		case "list":
			err = cli.User.ListUsers(conf, client)
		case "get":
			err = cli.User.GetUser(conf, client)
		case "update":
			err = cli.User.UpdateUser(conf, client)
		case "delete":
			err = cli.User.DeleteUser(conf, client)
		default:
			panic(fmt.Sprintf("unknown user command: %s", subcmd))
		}
	default:
		panic(fmt.Sprintf("unknown command: %s", kongCtx.Command()))

	}
	kongCtx.FatalIfErrorf(err)
}

// createAuthenticatedContext loads token using TokenManager and creates an authenticated gRPC context
func createAuthenticatedContext(jwtManager *authutils.JWTManager) (context.Context, error) {
	tokenManager := auth.NewTokenManager(jwtManager)
	ctx, err := tokenManager.CreateAuthenticatedContext(context.Background())
	if err != nil {
		// Provide a more user-friendly error message
		if err.Error() == "invalid key size" || err.Error() == "key not found" {
			return nil, fmt.Errorf("authentication required - please run 'featurectl auth login <email>' first")
		}
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	return ctx, nil
}
