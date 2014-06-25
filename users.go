package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"os"
)

type UserCreationRequest struct {
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
}

func CreateUser(loginname string, password string) {
	path := fmt.Sprintf("/apps/%s/users", globalConfig.AppId)
	headers := globalConfig.HttpHeaders("application/json")
	req := &UserCreationRequest{loginname, password}
	res := HttpPostJson(path, headers, req)
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}

var UsersCommands = []cli.Command{
	{
		Name:        "users:create",
		Usage:       "Create user",
		Description: `arguments: <loginname> <password>`,
		Action: func(c *cli.Context) {
			ShowCommandHelp(2, c)
			CreateUser(c.Args()[0], c.Args()[1])
		},
	},
}
