package main

import (
	"os"

	"github.com/codegangsta/cli"
)

const (
	ExitGeneralReason       = 1
	ExitIllegalNumberOfArgs = 2
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
	app.Version = "0.0.3"
	app.Commands = Flatten([][]cli.Command{
		LoginCommands,
		LogCommands,
		ServerCodeCommands,
		UsersCommands,
		WSEchoCommands,
	})
	setupFlags(app)
	app.Run(os.Args)
}

func Flatten(a [][]cli.Command) []cli.Command {
	b := []cli.Command{}
	for _, v := range a {
		for _, i := range v {
			b = append(b, i)
		}
	}
	return b
}
