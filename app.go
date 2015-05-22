package kiicli

import (
	"fmt"
	"github.com/tmtk75/cli"
	"strings"
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
		Usage: "set an app param",
		Args:  "<paramname> <paramvalue>",
		Action: func(c *cli.Context) {
			name, _ := c.ArgFor("paramname")
			value, _ := c.ArgFor("paramvalue")
			SetAppParam(name, value)
		},
	},
}
