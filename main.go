package main

import (
	"fmt"
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
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:  "login",
			Usage: "Login as AppAdmin",
			Action: func(c *cli.Context) {
				err := LoginAsAppAdmin()
				if err != nil {
					panic(err)
				}
			},
		},
		{
			Name:  "login:info",
			Usage: "Print login info",
			Action: func(c *cli.Context) {
				res := &OAuth2Response{}
				res.Load()
				fmt.Println(res.AccessToken)
			},
		},
		{
			Name:        "users:create",
			Usage:       "Create user",
			Description: `arguments: <loginname> <password>`,
			Action: func(c *cli.Context) {
				ShowCommandHelp(2, c)
				CreateUser(c.Args()[0], c.Args()[1])
			},
		},
		{
			Name:  "servercode:list",
			Usage: "List versions of server code",
			Flags: []cli.Flag{
				cli.BoolFlag{"quite, q", "Print only versionID"},
				cli.BoolFlag{"active, a", "Print only active one"},
			},
			Action: func(c *cli.Context) {
				ListServerCode(c.Bool("quite"), c.Bool("active"))
			},
		},
		{
			Name:  "servercode:deploy",
			Usage: "Deploy a server code",
			Flags: []cli.Flag{
				cli.BoolFlag{"activate", "Activate after deploying"},
			},
			Action: func(c *cli.Context) {
				ShowCommandHelp(1, c)
				DeployServerCode(c.Args()[0], c.Bool("activate"))
			},
		},
		{
			Name:  "servercode:get",
			Usage: "Get specified server code",
			Action: func(c *cli.Context) {
				if len(c.Args()) > 1 {
					cli.ShowCommandHelp(c, c.Command.Name)
					os.Exit(ExitIllegalNumberOfArgs)
				}
				var ver string
				if len(c.Args()) == 1 {
					ver = c.Args()[0]
				} else {
					vers := ListVersions()
					for _, v := range vers.Versions {
						if v.Active {
							ver = v.VersionId
							break
						}
					}
				}
				GetServerCode(ver)
			},
		},
		{
			Name:        "servercode:invoke",
			Usage:       "Invoke an entry point of server code",
			Description: "arguments: <entry-name> [version]",
			Action: func(c *cli.Context) {
				if len(c.Args()) > 2 || len(c.Args()) == 0 {
					cli.ShowCommandHelp(c, c.Command.Name)
					os.Exit(ExitIllegalNumberOfArgs)
				}
				version := "current"
				if len(c.Args()) == 2 {
					version = c.Args()[1]
				}
				InvokeServerCode(c.Args()[0], version)
			},
		},
		{
			Name:  "servercode:activate",
			Usage: "Activate a version",
			Action: func(c *cli.Context) {
				ShowCommandHelp(1, c)
				ActivateServerCode(c.Args()[0])
			},
		},
		{
			Name:  "servercode:delete",
			Usage: "Delete an entry point of server code",
			Action: func(c *cli.Context) {
				ShowCommandHelp(1, c)
				DeleteServerCode(c.Args()[0])
			},
		},
		{
			Name:  "log",
			Usage: "Print logs",
			Action: func(c *cli.Context) {
				StartLogging()
			},
		},
		{
			Name:  "server",
			Usage: "WebSocket echo server for testing",
			Action: func(c *cli.Context) {
				StartWSEchoServer()
			},
		},
	}
	setupFlags(app)
	app.Run(os.Args)
}
