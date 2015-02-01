package kiicli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/tmtk75/cli"
	goext "github.com/tmtk75/go-ext"
)

type Headers map[string]string

func DeployServerCode(serverCodePath string, activate bool) string {
	code, err := ioutil.ReadFile(serverCodePath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	p := Profile()
	path := fmt.Sprintf("/apps/%s/server-code", p.AppId)
	headers := p.HttpHeadersWithAuthorization("application/javascript")
	b := HttpPost(path, headers, bytes.NewReader(code)).Bytes()
	var ver map[string]string
	json.Unmarshal(b, &ver)
	fmt.Printf("versionID: %s\n", ver["versionID"])
	if activate {
		ActivateServerCode(ver["versionID"])
	}
	return ver["versionID"]
}

func InvokeServerCode(entryName string, version string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/server-code/versions/%s/%s", p.AppId, version, entryName)
	headers := p.HttpHeadersWithAuthorization("application/json")
	r := goext.OptionalReader(func() io.Reader { return strings.NewReader("{}") })
	b := HttpPost(path, headers, r).Bytes()
	fmt.Printf("%s\n", string(b))
}

type Versions struct {
	// Nested struct is convenient to unmarshal json string here.
	// The reason why it's not used is for sorting.
	Versions RawVersions `json:"versions"`
}

type RawVersions []RawVersion

func (self RawVersions) Len() int {
	return len(self)
}

func (self RawVersions) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self RawVersions) Less(i, j int) bool {
	return self[i].CreatedAt < self[j].CreatedAt
}

type RawVersion struct {
	VersionId  string `json:"versionID"`
	CreatedAt  int64  `json:"createdAt"`
	ModifiedAt int64  `json:"modifiedAt"`
	Active     bool   `json:"current"`
}

type Version struct {
	VersionId string
	CreatedAt time.Time
	Active    string
}

func (self *RawVersion) Version() *Version {
	a := "inactive"
	if self.Active {
		a = "active"
	}
	return &Version{
		VersionId: self.VersionId,
		CreatedAt: timeFromUnix(self.CreatedAt),
		Active:    a,
	}
}

func ListServerCode(quite bool, active bool) {
	vers := ListVersions()
	PrintVersions(vers, quite, active)
}

func ListVersions() *Versions {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/server-code/versions", p.AppId)
	headers := p.HttpHeadersWithAuthorization("")
	b := HttpGet(path, headers).Bytes()
	vers := Versions{}
	err := json.Unmarshal(b, &vers)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &vers
}

func PrintVersions(vers *Versions, quite bool, active bool) {
	sort.Sort(vers.Versions)
	for _, raw := range vers.Versions {
		v := raw.Version()
		t := v.CreatedAt.Format(time.RFC3339)
		if active && !raw.Active {
			continue
		}
		if quite {
			fmt.Printf("%s\n", v.VersionId)
		} else {
			fmt.Printf("%s\t%s\t%s\n", v.VersionId, t, v.Active)
		}
	}
}

func GetServerCode(version string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/server-code/versions/%s", p.AppId, version)
	headers := p.HttpHeadersWithAuthorization("")
	b := HttpGet(path, headers).Bytes()
	fmt.Printf("%s\n", string(b))
}

func ActivateServerCode(version string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/server-code/versions/current", p.AppId)
	headers := p.HttpHeadersWithAuthorization("text/plain")
	HttpPut(path, headers, strings.NewReader(version)).Bytes()
}

func DeleteServerCode(version string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/server-code/versions/%s", p.AppId, version)
	headers := p.HttpHeadersWithAuthorization("")
	b := HttpDelete(path, headers).Bytes()
	fmt.Printf("%s\n", string(b))
}

func AttachHookConfig(hookConfigPath, version string) {
	code, err := ioutil.ReadFile(hookConfigPath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	p := Profile()
	path := fmt.Sprintf("/apps/%s/hooks/versions/%s", p.AppId, version)
	headers := p.HttpHeadersWithAuthorization("application/vnd.kii.HooksDeploymentRequest+json")
	b := HttpPut(path, headers, bytes.NewReader(code)).Bytes()
	var ver map[string]interface{}
	json.Unmarshal(b, &ver)
}

func GetHookConfig(version string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/hooks/versions/%s", p.AppId, version)
	headers := p.HttpHeadersWithAuthorization("")
	b := HttpGet(path, headers).Bytes()
	fmt.Printf("%s\n", string(b))
}

func DeleteHookConfig(version string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/hooks/versions/%s", p.AppId, version)
	headers := p.HttpHeadersWithAuthorization("")
	b := HttpDelete(path, headers).Bytes()
	fmt.Printf("%s\n", string(b))
}

type ExecutionResult struct {
	Description string `json:"queryDescription"`
	Results     []struct {
		Id         string `json:"scheduleExecutionID"`
		Status     string `json:"status"`
		Name       string `json:"name"`
		StartedAt  int64  `json:"startedAt"`
		FinishedAt int64  `json:"finishedAt"`
	} `json:"results"`
}

func ListExecutions() {
	now := time.Now().Unix() * 1000
	dayBefore1week := now - 60*60*24*7*1000 // in millisecond
	p := Profile()
	path := fmt.Sprintf("/apps/%s/hooks/executions/query", p.AppId)
	headers := p.HttpHeadersWithAuthorization("application/vnd.kii.ScheduleExecutionQueryRequest+json")
	q := fmt.Sprintf(`{
		             "scheduleExecutionQuery": {
		               "clause": {
		                 "type": "range",
			         "field": "startedAt",
			         "lowerLimit": %v,
			         "upperLimit": %v,
			         "lowerIncluded": true,
			         "upperIncluded": true
		               },
			       "orderBy": "startedAt",
			       "descending": false
		             }
		           }`, dayBefore1week, now)
	b := HttpPost(path, headers, bytes.NewReader([]byte(q))).Bytes()
	var r ExecutionResult
	if err := json.Unmarshal(b, &r); err != nil {
		log.Fatalf("%v", err)
	}
	for _, e := range r.Results {
		s := timeFromUnix(e.StartedAt).Format(time.RFC3339)
		f := timeFromUnix(e.FinishedAt).Format(time.RFC3339)
		fmt.Printf("%v\t%v\t%v\t%v\t%v\n", e.Id, s, f, e.Status, e.Name)
	}
}

func getActiveVersion(c *cli.Context) string {
	ver, b := c.ArgFor("version")
	if b {
		return ver
	}
	vers := ListVersions()
	for _, v := range vers.Versions {
		if v.Active {
			return v.VersionId
		}
	}
	log.Fatalf("Missing active version")
	return ""
}

var ServerCodeCommands = []cli.Command{
	{
		Name:  "list",
		Usage: "List versions of server code",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "quite, q", Usage: "Print only versionID"},
			cli.BoolFlag{Name: "active, a", Usage: "Print only active one"},
		},
		Action: func(c *cli.Context) {
			ListServerCode(c.Bool("quite"), c.Bool("active"))
		},
	},
	{
		Name:  "deploy",
		Usage: "Deploy a server code",
		Args:  "<servercode-path>",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "activate,a", Usage: "Activate after deploying"},
			cli.StringFlag{Name: "config-file", Usage: "File path to a hook config"},
		},
		Action: func(c *cli.Context) {
			p, _ := c.ArgFor("servercode-path")
			version := DeployServerCode(p, c.Bool("activate"))
			if path := c.String("config-file"); path != "" {
				AttachHookConfig(path, version)
			}
		},
	},
	{
		Name:  "get",
		Usage: "Get specified server code",
		Args:  "[version]",
		Action: func(c *cli.Context) {
			ver := getActiveVersion(c)
			GetServerCode(ver)
		},
	},
	{
		Name:  "invoke",
		Usage: "Invoke an entry point of server code",
		Args:  "<entry-name> [version]",
		Action: func(c *cli.Context) {
			name, _ := c.ArgFor("entry-name")
			ver, b := c.ArgFor("version")
			if !b {
				ver = "current"
			}
			InvokeServerCode(name, ver)
		},
	},
	{
		Name:  "activate",
		Usage: "Activate a version",
		Args:  "<version>",
		Action: func(c *cli.Context) {
			ver, _ := c.ArgFor("version")
			ActivateServerCode(ver)
		},
	},
	{
		Name:  "delete",
		Usage: "Delete a version of server code",
		Args:  "<version>",
		Action: func(c *cli.Context) {
			ver, _ := c.ArgFor("version")
			DeleteServerCode(ver)
		},
	},
	{
		Name:  "hook-attach",
		Usage: "Attach a hook config to current or specified server code",
		Args:  "<hook-config-path> [version]",
		Description: `An example of definition of server hook config 

      {
        "kiicloud://users" : [ {
          "what" : "EXECUTE_SERVER_CODE",
          "when" : "USER_CREATED",
          "endpoint" : "main"
        } ],
        "kiicloud://scheduler" : {
          "HourlyMessage" : {
            "what" : "EXECUTE_SERVER_CODE",
            "name" : "HourlyMessage",
            "cron" : "15 * * * *",
            "endpoint" : "main",
            "parameters" : {"message" : "Hello"}
          }
        }
      }`,
		Action: func(c *cli.Context) {
			path, _ := c.ArgFor("hook-config-path")
			ver := getActiveVersion(c)
			AttachHookConfig(path, ver)
		},
	},
	{
		Name:  "hook-get",
		Usage: "Get hook the config of current or specified server code",
		Args:  "[version]",
		Action: func(c *cli.Context) {
			ver := getActiveVersion(c)
			GetHookConfig(ver)
		},
	},
	{
		Name:  "hook-delete",
		Usage: "Delete the hook config of current specified server code",
		Args:  "[version]",
		Action: func(c *cli.Context) {
			ver := getActiveVersion(c)
			DeleteHookConfig(ver)
		},
	},
	{
		Name:  "list-executions",
		Usage: "List executions for 7 days before",
		Action: func(c *cli.Context) {
			ListExecutions()
		},
	},
}
