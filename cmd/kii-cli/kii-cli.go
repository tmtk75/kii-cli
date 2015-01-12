package main

import (
	"os"

	"github.com/tmtk75/cli"
	kiicli "github.com/tmtk75/kii-cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "kii-cli"
	app.Usage = "KiiCloud command line interface"
	app.Version = "0.1.4"
	app.Commands = []cli.Command{
		cli.Command{
			Name:        "auth",
			Usage:       "Authentication",
			Subcommands: kiicli.LoginCommands,
		},
		cli.Command{
			Name:        "app",
			Usage:       "Application management",
			Subcommands: kiicli.AppCommands,
		},
		kiicli.LogCommands[0],
		cli.Command{
			Name:        "servercode",
			Usage:       "Server code management",
			Subcommands: kiicli.ServerCodeCommands,
		},
		cli.Command{
			Name:        "user",
			Usage:       "User management",
			Subcommands: kiicli.UserCommands,
		},
		cli.Command{
			Name:        "bucket",
			Usage:       "Bucket management",
			Subcommands: kiicli.BucketCommands,
		},
		cli.Command{
			Name:        "object",
			Usage:       "Object management",
			Subcommands: kiicli.ObjectCommands,
		},
		cli.Command{
			Name:        "dev",
			Usage:       "Development support",
			Subcommands: kiicli.WSEchoCommands,
		},
		cli.Command{
			Name:        "profile",
			Usage:       "Profile management",
			Subcommands: kiicli.ProfileCommands,
		},
	}
	if os.Getenv("FLAT") != "" {
		app.Commands = kiicli.Flatten(app.Commands)
	}
	kiicli.SetupFlags(app)
	app.Run(os.Args)
}
