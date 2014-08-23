package main

import (
	"os"

	"github.com/codegangsta/cli"
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

func main() {
	app := cli.NewApp()
	app.Name = "kii-cli"
	app.Usage = "KiiCloud command line interface"
	app.Version = "0.0.6"
	app.Commands = Flatten([][]cli.Command{
		LoginCommands,
		LogCommands,
		ServerCodeCommands,
		BucketCommands,
		UserCommands,
		ObjectCommands,
		WSEchoCommands,
	})
	setupFlags(app)
	app.Run(os.Args)
}

func countAll(a [][]cli.Command) int {
	c := 0
	for _, v := range a {
		c += len(v)
	}
	return c
}

func Flatten(a [][]cli.Command) []cli.Command {
	b := make([]cli.Command, 0, countAll(a))
	for _, v := range a {
		for _, i := range v {
			b = append(b, i)
		}
	}
	return b
}
