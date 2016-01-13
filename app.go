package kiicli

import (
	"fmt"
	"strings"

	"github.com/tmtk75/cli"
)

func PrintAppInfo() {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/", p.AppId)
	headers := p.HttpHeadersWithAuthorization("")
	body := HttpGet(path, headers).Bytes()
	fmt.Printf("%s\n", string(body))
}

func SetAppParam(name string, value string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/configuration/parameters/%s", p.AppId, name)
	headers := p.HttpHeadersWithAuthorization("text/plain")
	body := HttpPut(path, headers, strings.NewReader(value)).Bytes()
	fmt.Printf("%s\n", string(body))
}

func DeleteAppParam(prop string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/configuration/parameters/%s", p.AppId, prop)
	headers := p.HttpHeadersWithAuthorization("application/json")
	body := HttpDelete(path, headers).Bytes()
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
	{
		Name:  "set-param",
		Usage: "Set an app param",
		Args:  "<paramname> <paramvalue>",
		Action: func(c *cli.Context) {
			name, _ := c.ArgFor("paramname")
			value, _ := c.ArgFor("paramvalue")
			SetAppParam(name, value)
		},
	},
	{
		Name:  "delete-param",
		Usage: "Delete an app param",
		Args:  "<paramname>",
		Action: func(c *cli.Context) {
			name, _ := c.ArgFor("paramname")
			DeleteAppParam(name)
		},
	},
}
