package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/tmtk75/cli"
)

type OAuth2Request struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (self *GlobalConfig) OAuth2Request() *OAuth2Request {
	return &OAuth2Request{
		self.ClientId,
		self.ClientSecret,
	}
}

type OAuth2Response struct {
	Id          string `json:"id"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (self *OAuth2Response) Bytes() []byte {
	it, err := json.Marshal(self)
	if err != nil {
		panic(err)
	}
	return it
}

func (self *OAuth2Response) Save(filename string) {
	err := ioutil.WriteFile(filename, self.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

func (self *OAuth2Response) Load() *OAuth2Response {
	path := adminTokenFilePath()
	b, _ := exists(path)
	if !b {
		print("You've not logged in, first `login`\n")
		os.Exit(ExitNotLoggedIn)
	}
	self.LoadFrom(path)
	return self
}

func (self *OAuth2Response) LoadFrom(path string) *OAuth2Response {
	logger.Printf("Load %s\n", path)

	file, _ := os.Open(path)
	body, err := ioutil.ReadAll(bufio.NewReader(file))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	json.Unmarshal(body, self)
	return self
}

func (self *OAuth2Response) Decode(res *HttpResponse) *OAuth2Response {
	d := json.NewDecoder(res.Body)
	err := d.Decode(&self)
	if err != nil {
		panic(err)
	}
	return self
}

func retrieveAppAdminAccessToken() *OAuth2Response {
	token := globalConfig.OAuth2Request()
	headers := globalConfig.HttpHeaders("application/json")
	res := HttpPostJson("/oauth2/token", headers, token)
	oauth2res := &OAuth2Response{}
	return oauth2res.Decode(res)
}

func adminTokenFilePath() string {
	return metaFilePath(path.Join(".", globalConfig.AppId), "token")
}

func LoginAsAppAdmin(force bool) {
	if b, _ := exists(adminTokenFilePath()); b && !force {
		print("Already logged in, use `--force` to login\n")
		os.Exit(0)
	}
	res := retrieveAppAdminAccessToken()
	res.Save(adminTokenFilePath())
}

var LoginCommands = []cli.Command{
	{
		Name:  "login",
		Usage: "Login as AppAdmin",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "force", Usage: "Force to login"},
		},
		Action: func(c *cli.Context) {
			LoginAsAppAdmin(c.Bool("force"))
		},
	},
	{
		Name:  "info",
		Usage: "Print login info",
		Action: func(c *cli.Context) {
			res := &OAuth2Response{}
			res.Load()
			fmt.Println(res.AccessToken)
		},
	},
}
