package main

import (
	"fmt"

	"github.com/tmtk75/cli"
)

func PrintAppInfo() {
	path := fmt.Sprintf("/apps/%s/", globalConfig.AppId)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	body := HttpGet(path, headers).Bytes()
	fmt.Printf("%s\n", string(body))
}

var AppCommands = []cli.Command{
	{
		Name:  "app:config",
		Usage: "Print config of app",
		Action: func(c *cli.Context) {
			PrintAppInfo()
		},
	},
}
