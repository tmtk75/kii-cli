package kiicli

import (
	"fmt"

	"github.com/tmtk75/cli"
)

func ListGroups() {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/groups?is_member=%s", p.AppId)
	headers := p.HttpHeadersWithAuthorization("")
	//body := strings.NewReader(`{"userQuery":{"clause":{"type":"all"}}}`)
	b := HttpGet(path, headers).Bytes()
	fmt.Println(string(b))
}

var GroupCommands = []cli.Command{
	{
		Name:  "list",
		Usage: "List groups",
		Action: func(c *cli.Context) {
			ListGroups()
		},
	},
}
