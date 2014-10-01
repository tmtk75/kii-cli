package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/tmtk75/cli"
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

func LoginAsUser(username string, password string) {
	dir := path.Join(globalConfig.AppId, username)
	tokenPath := metaFilePath(dir, "token")

	oauth2res := &OAuth2Response{}
	if b, _ := exists(tokenPath); b {
		oauth2res.LoadFrom(tokenPath)
	} else {
		headers := globalConfig.HttpHeaders("application/json")
		req := map[string]string{"username": username, "password": password}
		res := HttpPostJson("/oauth2/token", headers, req)
		oauth2res.Decode(res)
		oauth2res.Save(tokenPath)
	}

	fmt.Println(oauth2res)
}

func ReadUser(userId string) {
	path := fmt.Sprintf("/apps/%s/users/%v", globalConfig.AppId, userId)
	headers := globalConfig.HttpHeadersWithAuthorization("application/json")
	b := HttpGet(path, headers).Bytes()
	fmt.Println(string(b))
}

func ListUsers() {
	path := fmt.Sprintf("/apps/%s/users/query", globalConfig.AppId)
	headers := globalConfig.HttpHeadersWithAuthorization("application/vnd.kii.userqueryrequest+json")
	body := strings.NewReader(`{"userQuery":{"clause":{"type":"all"}}}`)
	b := HttpPost(path, headers, body).Bytes()
	fmt.Println(string(b))
}

func DeleteUser(userId string) {
	path := fmt.Sprintf("/apps/%s/users/%v", globalConfig.AppId, userId)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	b := HttpDelete(path, headers).Bytes()
	fmt.Println(string(b))
}

var UserCommands = []cli.Command{
	{
		Name:        "user:create",
		Usage:       "Create a user",
		Description: `args: <loginname> <password>`,
		Action: func(c *cli.Context) {
			ShowCommandHelp(2, c)
			CreateUser(c.Args()[0], c.Args()[1])
		},
	},
	{
		Name:        "user:read",
		Usage:       "Read a user",
		Description: `args: <userId>`,
		Action: func(c *cli.Context) {
			ShowCommandHelp(1, c)
			ReadUser(c.Args()[0])
		},
	},
	{
		Name:  "user:list",
		Usage: "List users",
		Action: func(c *cli.Context) {
			ShowCommandHelp(0, c)
			ListUsers()
		},
	},
	{
		Name:        "user:delete",
		Usage:       "Delete a user",
		Description: `args: <userId>`,
		Action: func(c *cli.Context) {
			ShowCommandHelp(1, c)
			DeleteUser(c.Args()[0])
		},
	},
	{
		Name:        "user:login",
		Usage:       "Login as a user",
		Description: `args: <loginname> <password>`,
		Action: func(c *cli.Context) {
			ShowCommandHelp(2, c)
			LoginAsUser(c.Args()[0], c.Args()[1])
		},
	},
}
