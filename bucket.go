package kiicli

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/tmtk75/cli"
)

func ListBucket() {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/buckets", p.AppId)
	headers := p.HttpHeadersWithAuthorization("")
	body := HttpGet(path, headers).Bytes()

	var v map[string][]interface{}
	json.Unmarshal(body, &v)
	for _, a := range v["bucketIDs"] {
		fmt.Println(a)
	}
}

func ReadBucket(bucketname string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/buckets/%s", p.AppId, bucketname)
	headers := p.HttpHeadersWithAuthorization("")
	body := HttpGet(path, headers).Bytes()
	fmt.Println(string(body))
}

func DeleteBucket(name string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/buckets/%v", p.AppId, name)
	headers := p.HttpHeadersWithAuthorization("")
	HttpDelete(path, headers).Bytes()
}

func readBucketAcl(bucketname string) []byte {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/buckets/%s/acl", p.AppId, bucketname)
	headers := p.HttpHeadersWithAuthorization("")
	body := HttpGet(path, headers).Bytes()
	return body
}

func ReadBucketAcl(bucketname string) {
	body := readBucketAcl(bucketname)
	fmt.Println(string(body))
}

func DeleteBucketAcl(bucketname, verb, userId string) {
	p := Profile()
	// /apps/%s/buckets/%s/acl/QUERY_OBJECTS_IN_BUCKET/UserID:ANONYMOUS_USER
	path := fmt.Sprintf("/apps/%s/buckets/%s/acl/%v/UserID:%v", p.AppId, bucketname, verb, userId)
	headers := p.HttpHeadersWithAuthorization("")
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
			//logger.Printf("%v:%v\n", verb, e.UserId)
			DeleteBucketAcl(bucketname, verb, e.UserId)
		}
	}
}

var BucketCommands = []cli.Command{
	{
		Name:  "list",
		Usage: "List buckets",
		Action: func(c *cli.Context) {
			ListBucket()
		},
	},
	{
		Name:  "read",
		Usage: "Read a bucket",
		Args:  "<bucket-id>",
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			ReadBucket(bid)
		},
	},
	{
		Name:  "delete",
		Usage: "Delete a bucket",
		Args:  "<bucket-id>",
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			DeleteBucket(bid)
		},
	},
	{
		Name:  "acl",
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
		Name:  "read",
		Usage: "Read a bucket ACL",
		Args:  "<bucket-id>",
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			ReadBucketAcl(bid)
		},
	},
	{
		Name:  "delete",
		Usage: "Delete a bucket ACL",
		Args:  `<bucket-id> <verb> <user-id>`,
		Description: `\
   ex)  my-bucket CREATE_OBJECTS_IN_BUCKET ANONYMOUS_USER
        my-bucket QUERY_OBJECTS_IN_BUCKET ANY_AUTHENTICATED_USER`,
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			verb, _ := c.ArgFor("verb")
			uid, _ := c.ArgFor("user-id")
			DeleteBucketAcl(bid, verb, uid)
		},
	},
	{
		Name:  "delete-all",
		Usage: "Delete all ACLs",
		Args:  "<bucket-id>",
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			DeleteAllBucketAcls(bid)
		},
	},
}
