package main

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/codegangsta/cli"
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
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

func LoginAsUser(username string, password string) {
	headers := globalConfig.HttpHeaders("application/json")
	req := map[string]string{"username": username, "password": password}
	res := HttpPostJson("/oauth2/token", headers, req)
	oauth2res := &OAuth2Response{}
	oauth2res.Decode(res)

	dir := path.Join(globalConfig.AppId, oauth2res.Id)
	oauth2res.Save(metaFilePath(dir, username))

	fmt.Println(oauth2res)
}

var UserCommands = []cli.Command{
	{
		Name:        "user:login",
		Usage:       "Login as a user",
		Description: `arguments: <loginname> <password>`,
		Action: func(c *cli.Context) {
			ShowCommandHelp(2, c)
			LoginAsUser(c.Args()[0], c.Args()[1])
		},
	},
	{
		Name:        "user:create",
		Usage:       "Create a user",
		Description: `arguments: <loginname> <password>`,
		Action: func(c *cli.Context) {
			ShowCommandHelp(2, c)
			CreateUser(c.Args()[0], c.Args()[1])
		},
	},
}
