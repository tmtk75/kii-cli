package kiicli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/mitchellh/colorstring"
	"github.com/tmtk75/cli"
	goext "github.com/tmtk75/go-ext"
)

func createObject(bucketname string, r io.Reader) map[string]interface{} {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/buckets/%s/objects", p.AppId, bucketname)
	headers := p.HttpHeadersWithAuthorization("application/json")
	body := HttpPost(path, headers, r).Bytes()
	var j map[string]interface{}
	json.Unmarshal(body, &j)
	return j
}

func CreateObject(bucketname string) {
	r := goext.OptionalReader(func() io.Reader { return strings.NewReader("{}") })
	j := createObject(bucketname, r)
	fmt.Printf("%s\n", j["objectID"])
}

func ReadObject(bucketname, objectId, templstr string) {
	var templ *template.Template
	if templstr != "" {
		t, err := template.New("").Parse(templstr)
		if err != nil {
			log.Fatalf("%v", err)
		}
		templ = t
	}

	p := Profile()
	path := fmt.Sprintf("/apps/%s/buckets/%s/objects/%s", p.AppId, bucketname, objectId)
	headers := p.HttpHeadersWithAuthorization("")
	body := HttpGet(path, headers).Bytes()

	if templ == nil {
		fmt.Printf("%s\n", string(body))
	} else {
		var j map[string]interface{}
		json.Unmarshal(body, &j)
		templ.Execute(os.Stdout, j)
	}
}

func QueryObject(bucketname string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/buckets/%s/query", p.AppId, bucketname)
	headers := p.HttpHeadersWithAuthorization("application/vnd.kii.QueryRequest+json")
	r := goext.OptionalReader(func() io.Reader { return strings.NewReader(`{"bucketQuery":{"clause":{"type":"all"}}}`) })
	body := HttpPost(path, headers, r).Bytes()

	fmt.Printf("%s\n", string(body))
}

func ReplaceObject(bucketname string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/buckets/%s/objects", p.AppId, bucketname)
	headers := p.HttpHeadersWithAuthorization("application/json")
	r := goext.OptionalReader(func() io.Reader { return strings.NewReader("{}") })
	body := HttpPut(path, headers, r).Bytes()

	var j map[string]interface{}
	json.Unmarshal(body, &j)
	fmt.Printf("%s\n", j["objectID"])
}

func DeleteObject(bucketname, objectId string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/buckets/%s/objects/%s", p.AppId, bucketname, objectId)
	headers := p.HttpHeadersWithAuthorization("")
	body := HttpDelete(path, headers).Bytes()
	fmt.Printf("%s\n", string(body))
}

func attachObjectBody(bucketname, objectId, conttype string, r io.Reader) []byte {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/buckets/%s/objects/%v/body", p.AppId, bucketname, objectId)
	headers := p.HttpHeadersWithAuthorization(conttype)
	return HttpPut(path, headers, r).Bytes()
}

func AttachObjectBody(bucketname, objectId, conttype string, r io.Reader) {
	body := attachObjectBody(bucketname, objectId, conttype, r)
	fmt.Printf("%v", string(body))
}

func publishObjectBody(bucketname, objectId string, d ExpiredDuration) []byte {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/buckets/%s/objects/%v/body/publish", p.AppId, bucketname, objectId)
	headers := p.HttpHeadersWithAuthorization("application/vnd.kii.ObjectBodyPublicationRequest+json")
	req := map[string]int64{"expiresIn": (int64(d))}
	j, _ := json.Marshal(req)
	return HttpPost(path, headers, bytes.NewReader(j)).Bytes()
}

func PublishObjectBody(bucketname, objectId string, d ExpiredDuration) {
	body := publishObjectBody(bucketname, objectId, d)
	fmt.Printf("%v", string(body))
}

func CreateObjectAndPublishBody(bucketname, conttype string, d ExpiredDuration, br io.Reader) {
	r := strings.NewReader("{}")
	j := createObject(bucketname, r)
	objId := j["objectID"].(string)

	r0 := attachObjectBody(bucketname, objId, conttype, br)
	var a map[string]int64
	json.Unmarshal(r0, &a)
	logger.Printf("modifiedAt: %v", a["modifiedAt"])

	r1 := publishObjectBody(bucketname, objId, d)
	var res struct {
		PublicationID string `json:"publicationID"`
		URL           string `json:"url"`
	}
	json.Unmarshal(r1, &res)
	logger.Printf("publicationID: %v", res.PublicationID)
	logger.Printf("url: %v", res.URL)

	fmt.Printf("%v\n", res.URL)
}

var ObjectCommands = []cli.Command{
	{
		Name:  "create",
		Usage: "Create an object in application scope",
		Args:  "<bucket-id>",
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			CreateObject(bid)
		},
	},
	{
		Name:  "read",
		Usage: "Read the object in application scope",
		Args:  "<bucket-id> <object-id>",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "template", Value: "", Usage: "Template for output"},
		},
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			oid, _ := c.ArgFor("object-id")
			ReadObject(bid, oid, c.String("template"))
		},
	},
	{
		Name:  "query",
		Usage: "Query objects in a bucket of application scope",
		Args:  "<bucket-id>",
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			QueryObject(bid)
		},
	},
	//{
	//	Name:        "replace",
	//	Usage:       "Replate the object in application scope with a new one",
	//	Description: "args: <bucket> <object-id>",
	//	Action: func(c *cli.Context) {
	//		ShowCommandHelp(2, c)
	//		ReplaceObject(c.Args()[0], c.Args()[1])
	//	},
	//},
	{
		Name:  "delete",
		Usage: "Delete the object in application scope",
		Args:  "<bucket-id> <object-id>",
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			oid, _ := c.ArgFor("object-id")
			DeleteObject(bid, oid)
		},
	},
	{
		Name:  "body-attach",
		Usage: "Attach body to an object in application scope",
		Args:  `<bucket-id> <object-id> <content-type>`,
		Description: `
   ex)
     dogs 4c8aaf60-3166-11e4-a448-12315004cc43 image/png < mydog.png
`,
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			oid, _ := c.ArgFor("object-id")
			ctype, _ := c.ArgFor("content-type")
			AttachObjectBody(bid, oid, ctype, os.Stdin)
		},
	},
	{
		Name:  "body-publish",
		Usage: "Publish a body of object in application scope",
		Args:  `<bucket-id> <object-id>`,
		Flags: publishFlags,
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			oid, _ := c.ArgFor("object-id")
			expired := relativeTimeToDuration(c.String("expired-in"))
			PublishObjectBody(bid, oid, expired)
		},
	},
	{
		Name:  "publish",
		Usage: "Publish a body creating a new object into the bucket in application scope",
		Args:  `<bucket-id> <content-type>`,
		Flags: publishFlags,
		Description: `
   Runs object:create, object-body-attach and object:body-publish in order.
   It's expected body is given thru stdin.

   ex)
     dogs image/png < mydog.png`,
		Action: func(c *cli.Context) {
			bid, _ := c.ArgFor("bucket-id")
			ctype, _ := c.ArgFor("content-type")
			expired := relativeTimeToDuration(c.String("expired-in"))
			r := goext.OptionalReader(func() io.Reader {
				log.Fatalf(colorstring.Color("[red]object body must be given thru stdin"))
				return nil
			})
			CreateObjectAndPublishBody(bid, ctype, expired, r)
		},
	},
}

type ExpiredDuration int64

var publishFlags = []cli.Flag{
	cli.StringFlag{Name: "expired-in", Value: "1m", Usage: "Duration in seconds the publication URL has to be available, after that it will expire"},
}

func relativeTimeToDuration(s string) ExpiredDuration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ExpiredDuration(d.Seconds())
}
