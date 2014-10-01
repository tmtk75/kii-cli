package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/tmtk75/cli"

	"code.google.com/p/go.crypto/ssh/terminal"
)

type Headers map[string]string

func DeployServerCode(serverCodePath string, activate bool) string {
	code, err := ioutil.ReadFile(serverCodePath)
	if err != nil {
		panic(err)
	}
	path := fmt.Sprintf("/apps/%s/server-code", globalConfig.AppId)
	headers := globalConfig.HttpHeadersWithAuthorization("application/javascript")
	b := HttpPost(path, headers, bytes.NewReader(code)).Bytes()
	var ver map[string]string
	json.Unmarshal(b, &ver)
	fmt.Printf("versionID: %s\n", ver["versionID"])
	if activate {
		ActivateServerCode(ver["versionID"])
	}
	return ver["versionID"]
}

func OptionalReader(f func() io.Reader) io.Reader {
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		return f()
	}
	return os.Stdin
}

func InvokeServerCode(entryName string, version string) {
	path := fmt.Sprintf("/apps/%s/server-code/versions/%s/%s", globalConfig.AppId, version, entryName)
	headers := globalConfig.HttpHeadersWithAuthorization("application/json")
	r := OptionalReader(func() io.Reader { return strings.NewReader("{}") })
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
	return &Version{self.VersionId, time.Unix(self.CreatedAt/1000, self.CreatedAt%1000*1000*1000), a}
}

func ListServerCode(quite bool, active bool) {
	vers := ListVersions()
	PrintVersions(vers, quite, active)
}

func ListVersions() *Versions {
	path := fmt.Sprintf("/apps/%s/server-code/versions", globalConfig.AppId)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	b := HttpGet(path, headers).Bytes()
	vers := Versions{}
	err := json.Unmarshal(b, &vers)
	if err != nil {
		panic(err)
	}
	return &vers
}

func PrintVersions(vers *Versions, quite bool, active bool) {
	sort.Sort(vers.Versions)
	for _, raw := range vers.Versions {
		v := raw.Version()
		t := v.CreatedAt.Format("2006-01-02 15:04:05")
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
	path := fmt.Sprintf("/apps/%s/server-code/versions/%s", globalConfig.AppId, version)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	b := HttpGet(path, headers).Bytes()
	fmt.Printf("%s\n", string(b))
}

func ActivateServerCode(version string) {
	path := fmt.Sprintf("/apps/%s/server-code/versions/current", globalConfig.AppId)
	headers := globalConfig.HttpHeadersWithAuthorization("text/plain")
	HttpPut(path, headers, strings.NewReader(version)).Bytes()
}

func DeleteServerCode(version string) {
	path := fmt.Sprintf("/apps/%s/server-code/versions/%s", globalConfig.AppId, version)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	b := HttpDelete(path, headers).Bytes()
	fmt.Printf("%s\n", string(b))
}

func AttachHookConfig(hookConfigPath, version string) {
	code, err := ioutil.ReadFile(hookConfigPath)
	if err != nil {
		panic(err)
	}
	path := fmt.Sprintf("/apps/%s/hooks/versions/%s", globalConfig.AppId, version)
	headers := globalConfig.HttpHeadersWithAuthorization("application/vnd.kii.HooksDeploymentRequest+json")
	b := HttpPut(path, headers, bytes.NewReader(code)).Bytes()
	var ver map[string]interface{}
	json.Unmarshal(b, &ver)
}

func GetHookConfig(version string) {
	path := fmt.Sprintf("/apps/%s/hooks/versions/%s", globalConfig.AppId, version)
	headers := globalConfig.HttpHeadersWithAuthorization("")
	b := HttpGet(path, headers).Bytes()
	fmt.Printf("%s\n", string(b))
}

func DeleteHookConfig(version string) {
	path := fmt.Sprintf("/apps/%s/hooks/versions/%s", globalConfig.AppId, version)
	headers := globalConfig.HttpHeadersWithAuthorization("")
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
	path := fmt.Sprintf("/apps/%s/hooks/executions/query", globalConfig.AppId)
	headers := globalConfig.HttpHeadersWithAuthorization("application/vnd.kii.ScheduleExecutionQueryRequest+json")
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
		panic(err)
	}
	for _, e := range r.Results {
		s := time.Unix(e.StartedAt/1000, 0).Format("2006-01-02 15:04:05")
		f := time.Unix(e.FinishedAt/1000, 0).Format("2006-01-02 15:04:05")
		fmt.Printf("%v\t%v\t%v\t%v\t%v\n", e.Id, s, f, e.Status, e.Name)
	}
}

func getActiveVersion(c *cli.Context, argLen int) string {
	if len(c.Args()) > argLen {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(ExitIllegalNumberOfArgs)
	}
	if len(c.Args()) == argLen {
		return c.Args()[argLen-1]
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
		Name:  "servercode:list",
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
		Name:  "servercode:deploy",
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
		Name:  "servercode:get",
		Usage: "Get specified server code",
		Args:  "<version>",
		Action: func(c *cli.Context) {
			ver, _ := c.ArgFor("version")
			GetServerCode(ver)
		},
	},
	{
		Name:  "servercode:invoke",
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
		Name:  "servercode:activate",
		Usage: "Activate a version",
		Args:  "<version>",
		Action: func(c *cli.Context) {
			ver, _ := c.ArgFor("version")
			ActivateServerCode(ver)
		},
	},
	{
		Name:  "servercode:delete",
		Usage: "Delete a version of server code",
		Args:  "<version>",
		Action: func(c *cli.Context) {
			ver, _ := c.ArgFor("version")
			DeleteServerCode(ver)
		},
	},
	{
		Name:  "servercode:hook-attach",
		Usage: "Attach a hook config to current or specified server code",
		Args:  "<hook-config-path> [version]",
		Action: func(c *cli.Context) {
			path, _ := c.ArgFor("hook-config-path")
			ver, b := c.ArgFor("version")
			if !b {
				ver = getActiveVersion(c, 2)
			}
			AttachHookConfig(path, ver)
		},
	},
	{
		Name:  "servercode:hook-get",
		Usage: "Get hook the config of current or specified server code",
		Args:  "[version]",
		Action: func(c *cli.Context) {
			ver, b := c.ArgFor("version")
			if !b {
				ver = getActiveVersion(c, 1)
			}
			GetHookConfig(ver)
		},
	},
	{
		Name:  "servercode:hook-delete",
		Usage: "Delete the hook config of current specified server code",
		Args:  "[version]",
		Action: func(c *cli.Context) {
			ver, b := c.ArgFor("version")
			if !b {
				ver = getActiveVersion(c, 1)
			}
			DeleteHookConfig(ver)
		},
	},
	{
		Name:  "servercode:list-executions",
		Usage: "List executions for 7 days before",
		Action: func(c *cli.Context) {
			ListExecutions()
		},
	},
}
