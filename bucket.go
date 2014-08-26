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

func DeleteBucket(name string) {
	path := fmt.Sprintf("/apps/%s/buckets/%v", globalConfig.AppId, name)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	HttpDelete(path, headers).Bytes()
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
		Name:  "bucket:delete",
		Usage: "Delete a bucket",
		Action: func(c *cli.Context) {
			ShowCommandHelp(1, c)
			DeleteBucket(c.Args()[0])
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
