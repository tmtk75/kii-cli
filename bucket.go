package main

import (
	"encoding/json"
	"fmt"
	"log"

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

func readBucketAcl(bucketname string) []byte {
	path := fmt.Sprintf("/apps/%s/buckets/%s/acl", globalConfig.AppId, bucketname)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	body := HttpGet(path, headers).Bytes()
	return body
}
func ReadBucketAcl(bucketname string) {
	body := readBucketAcl(bucketname)
	fmt.Println(string(body))
}

func DeleteBucketAcl(bucketname, verb, userId string) {
	// /apps/%s/buckets/%s/acl/QUERY_OBJECTS_IN_BUCKET/UserID:ANONYMOUS_USER
	path := fmt.Sprintf("/apps/%s/buckets/%s/acl/%v/UserID:%v", globalConfig.AppId, bucketname, verb, userId)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	body := HttpDelete(path, headers).Bytes()
	fmt.Println(string(body))
}

func DeleteAllBucketAcls(bucketname string) {
	body := readBucketAcl(bucketname)
	var j map[string][]struct {
		UserId string `json:"userID"`
	}
	err := json.Unmarshal(body, &j)
	if err != nil {
		log.Fatalf("%v", err)
	}
	for verb, v := range j {
		for _, e := range v {
			logger.Printf("%v:%v\n", verb, e.UserId)
			DeleteBucketAcl(bucketname, verb, e.UserId)
		}
	}
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
		Usage: "Edit bucket ACL",
		Description: `Edit bucket ACL

     verb:
         CREATE_OBJECTS_IN_BUCKET
         QUERY_OBJECTS_IN_BUCKET
         DROP_BUCKET_WITH_ALL_CONTENT

     Special userID:
         ANY_AUTHENTICATED_USER
	 ANONYMOUS_USER`,
		Subcommands: cmds,
	},
}

var cmds = []cli.Command{
	{
		Name:        "read",
		Usage:       "Read a bucket ACL",
		Description: "args: <bucket>",
		Action: func(c *cli.Context) {
			ShowCommandHelp(1, c)
			ReadBucketAcl(c.Args()[0])
		},
	},
	{
		Name:  "delete",
		Usage: "Delete a bucket ACL",
		Description: `args: <bucket> <verb> <userID>

   ex)  my-bucket CREATE_OBJECTS_IN_BUCKET ANONYMOUS_USER
        my-bucket QUERY_OBJECTS_IN_BUCKET ANY_AUTHENTICATED_USER`,
		Action: func(c *cli.Context) {
			ShowCommandHelp(3, c)
			DeleteBucketAcl(c.Args()[0], c.Args()[1], c.Args()[2])
		},
	},
	{
		Name:        "delete-all",
		Usage:       "Delete all ACLs",
		Description: "args: <bucket>",
		Action: func(c *cli.Context) {
			ShowCommandHelp(1, c)
			DeleteAllBucketAcls(c.Args()[0])
		},
	},
}
