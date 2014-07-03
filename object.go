package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/codegangsta/cli"
)

func CreateObject(bucketname string) {
	path := fmt.Sprintf("/apps/%s/buckets/%s/objects", globalConfig.AppId, bucketname)
	headers := globalConfig.HttpHeadersWithAuthorization("application/json")
	r := OptionalReader(func() io.Reader { return strings.NewReader("{}") })
	body := HttpPost(path, headers, r).Bytes()

	var j map[string]interface{}
	json.Unmarshal(body, &j)
	fmt.Printf("%s\n", j["objectID"])
}

func ReadObject(bucketname, objectId string) {
	path := fmt.Sprintf("/apps/%s/buckets/%s/objects/%s", globalConfig.AppId, bucketname, objectId)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	body := HttpGet(path, headers).Bytes()
	fmt.Printf("%s\n", string(body))
}

var ObjectCommands = []cli.Command{
	{
		Name:  "object:create",
		Usage: "Create an object",
		Action: func(c *cli.Context) {
			ShowCommandHelp(1, c)
			CreateObject(c.Args()[0])
		},
	},
	{
		Name:  "object:read",
		Usage: "Read the object",
		Action: func(c *cli.Context) {
			ShowCommandHelp(2, c)
			ReadObject(c.Args()[0], c.Args()[1])
		},
	},
}
