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

func SetValueInAppConfig(prop, val string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/configuration/parameters/%s", p.AppId, prop)
	headers := p.HttpHeadersWithAuthorization("text/plain")
	body := HttpPut(path, headers, strings.NewReader(val)).Bytes()
	fmt.Printf("%s\n", string(body))
}

func DeleteValueInAppConfig(prop string) {
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
		Name:        "config-set",
		Usage:       "Set a value in a property of config",
		Args:        "<property> <value>",
		Description: `   e.g) phoneNumberVerificationRequired true`,
		Action: func(c *cli.Context) {
			prop, _ := c.ArgFor("property")
			val, _ := c.ArgFor("value")
			SetValueInAppConfig(prop, val)
		},
	},
	{
		Name:        "config-delete",
		Usage:       "Delete a property of config",
		Args:        "<property>",
		Description: `   e.g) kii.consumer_key`,
		Action: func(c *cli.Context) {
			prop, _ := c.ArgFor("property")
			DeleteValueInAppConfig(prop)
		},
	},
}
