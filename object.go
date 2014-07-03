package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"
)

func CreateObject(bucketname, filename string) {
	path := fmt.Sprintf("/apps/%s/buckets/%s/objects", globalConfig.AppId, bucketname)
	headers := globalConfig.HttpHeadersWithAuthorization("application/json")
	b, _ := ioutil.ReadFile(filename)
	r := bytes.NewReader(b)
	body := HttpPost(path, headers, r).Bytes()

	fmt.Println(string(body))
}

var ObjectCommands = []cli.Command{
	{
		Name:  "object:create",
		Usage: "Create an object",
		Action: func(c *cli.Context) {
			ShowCommandHelp(2, c)
			CreateObject(c.Args()[0], c.Args()[1])
		},
	},
}
