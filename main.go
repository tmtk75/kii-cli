package main

import (

	//"reflect"

	//"log"

	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "kii-cli"
	app.Usage = "KiiCloud command line tools"
	app.Commands = []cli.Command{
		{
			Name:  "log",
			Usage: "print logs",
			Action: func(c *cli.Context) {
				StartLogging()
			},
		},
	}
	app.Run(os.Args)
}
