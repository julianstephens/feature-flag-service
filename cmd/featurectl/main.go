package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/julianstephens/feature-flag-service/internal/commands"
	"github.com/julianstephens/feature-flag-service/internal/config"
)

type Globals struct {
	Version kong.VersionFlag
}

type CLI struct {
	Globals

	Login struct {
	} `cmd:"" help:"Login to the feature management system."`
	Flag  commands.FlagCommand `cmd:"" help:"Manage feature flags."`
	Audit struct {
	} `cmd:"" help:"Audit log operations."`
}

var err error

func main() {
	cli := CLI{
		Globals: Globals{},
	}

	conf := config.LoadConfig()

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
		kong.ConfigureHelp(kong.HelpOptions{Compact: true}),
		kong.Vars{"version": "1.0.0"},
	)

	cmd := strings.Split(kongCtx.Command(), " ")
	switch cmd[0] {
	case "login":
		// Implement login functionality here
	case "flag":
		if len(cmd) < 2 {
			kongCtx.PrintUsage(false)
			return
		}
		subcmd := cmd[1]
		switch subcmd {
		case "list":
			err = cli.Flag.ListFlags(conf, conn)
		case "get":
			err = cli.Flag.GetFlag(conf, conn)
		case "create":
			err = cli.Flag.CreateFlag(conf, conn)
		case "update":
			err = cli.Flag.UpdateFlag(conf, conn)
		case "delete":
			err = cli.Flag.DeleteFlag(conf, conn)
		default:
			panic(fmt.Sprintf("unknown flag command: %s", subcmd))
		}
	case "audit":
		// Implement audit log functionality here
	default:
		panic(fmt.Sprintf("unknown command: %s", kongCtx.Command()))

	}
	kongCtx.FatalIfErrorf(err)
}
