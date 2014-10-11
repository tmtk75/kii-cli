package main

import (
	"os"

	"github.com/tmtk75/cli"
)

const (
	ExitGeneralReason       = 1
	ExitIllegalNumberOfArgs = 2
	ExitNotLoggedIn         = 3
	ExitMissingParams       = 4
)

func ShowCommandHelp(argsLen int, c *cli.Context) {
	if len(c.Args()) != argsLen {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(ExitIllegalNumberOfArgs)
	}
}

type pair struct {
	prefix   cli.Command
	commands []cli.Command
}

var commands = []pair{}

func main() {
	app := cli.NewApp()
	app.Name = "kii-cli"
	app.Usage = "KiiCloud command line interface"
	app.Version = "0.1.2"
	app.Commands = []cli.Command{
		cli.Command{
			Name:        "auth",
			Usage:       "Authentication",
			Subcommands: LoginCommands,
		},
		cli.Command{
			Name:        "app",
			Usage:       "Application management",
			Subcommands: AppCommands,
		},
		LogCommands[0],
		cli.Command{
			Name:        "servercode",
			Usage:       "Server code management",
			Subcommands: ServerCodeCommands,
		},
		cli.Command{
			Name:        "user",
			Usage:       "User management",
			Subcommands: UserCommands,
		},
		cli.Command{
			Name:        "bucket",
			Usage:       "Bucket management",
			Subcommands: BucketCommands,
		},
		cli.Command{
			Name:        "object",
			Usage:       "Object management",
			Subcommands: ObjectCommands,
		},
		cli.Command{
			Name:        "dev",
			Usage:       "Development support",
			Subcommands: WSEchoCommands,
		},
	}
	if os.Getenv("FLAT") != "" {
		app.Commands = Flatten(app.Commands)
	}
	SetupFlags(app)
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
