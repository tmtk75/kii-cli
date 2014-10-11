package main

import (
	"fmt"

	"github.com/tmtk75/cli"
)

func PrintAppInfo() {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/", p.AppId)
	headers := p.HttpHeadersWithAuthorization("")
	body := HttpGet(path, headers).Bytes()
	fmt.Printf("%s\n", string(body))
}

var AppCommands = []cli.Command{
	{
		Name:  "config",
		Usage: "Print config of app",
		Action: func(c *cli.Context) {
			PrintAppInfo()
		},
	},
}
