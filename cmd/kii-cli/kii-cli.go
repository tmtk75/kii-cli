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
	app.Version = "0.1.3"
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
	}
	if os.Getenv("FLAT") != "" {
		app.Commands = Flatten(app.Commands)
	}
	kiicli.SetupFlags(app)
	app.Run(os.Args)
}

func Flatten(a []cli.Command) []cli.Command {
	b := make([]cli.Command, 0, 16)
	for _, v := range a {
		if v.Subcommands == nil {
			b = append(b, v)
		} else {
			for _, i := range v.Subcommands {
				i.Name = v.Name + ":" + i.Name
				b = append(b, i)
			}
		}
	}
	return b
}
