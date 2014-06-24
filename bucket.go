package main

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
)

func ListBucket() {
	path := fmt.Sprintf("/apps/%s/buckets", globalConfig.AppId)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	body := HttpGet(path, headers).Bytes()

	var v map[string][]interface{}
	json.Unmarshal(body, &v)
	for _, a := range v["bucketIDs"] {
		fmt.Println(a)
	}
}

func ShowBucketAcl(bucketname string) {
	path := fmt.Sprintf("/apps/%s/buckets/%s/acl", globalConfig.AppId, bucketname)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	body := HttpGet(path, headers).Bytes()
	fmt.Println(string(body))
}

var BucketCommands = []cli.Command{
	{
		Name:  "bucket:list",
		Usage: "List buckets",
		Action: func(c *cli.Context) {
			ListBucket()
		},
	},
	{
		Name:  "bucket:acl",
		Usage: "Show a bucket ACL",
		Action: func(c *cli.Context) {
			ShowCommandHelp(1, c)
			ShowBucketAcl(c.Args()[0])
		},
	},
}
