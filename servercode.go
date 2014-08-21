package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/codegangsta/cli"

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

func DeployHookConfig(hookConfigPath, version string) {
	code, err := ioutil.ReadFile(hookConfigPath)
	if err != nil {
		panic(err)
	}
	path := fmt.Sprintf("/apps/%s/hooks/versions/%s", globalConfig.AppId, version)
	headers := globalConfig.HttpHeadersWithAuthorization("application/vnd.kii.HooksDeploymentRequest+json")
	b := HttpPut(path, headers, bytes.NewReader(code)).Bytes()
	var ver map[string]interface{}
	json.Unmarshal(b, &ver)
	fmt.Printf("%v", ver)
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
		Name:        "servercode:deploy",
		Usage:       "Deploy a server code",
		Description: "args: <servercode-path>",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "activate", Usage: "Activate after deploying"},
			cli.StringFlag{Name: "hook-config-path", Usage: "File path to a hook config"},
		},
		Action: func(c *cli.Context) {
			ShowCommandHelp(1, c)
			version := DeployServerCode(c.Args()[0], c.Bool("activate"))
			if path := c.String("hook-config-path"); path != "" {
				DeployHookConfig(path, version)
			}
		},
	},
	{
		Name:  "servercode:get",
		Usage: "Get specified server code",
		Action: func(c *cli.Context) {
			if len(c.Args()) > 1 {
				cli.ShowCommandHelp(c, c.Command.Name)
				os.Exit(ExitIllegalNumberOfArgs)
			}
			var ver string
			if len(c.Args()) == 1 {
				ver = c.Args()[0]
			} else {
				vers := ListVersions()
				for _, v := range vers.Versions {
					if v.Active {
						ver = v.VersionId
						break
					}
				}
			}
			GetServerCode(ver)
		},
	},
	{
		Name:        "servercode:invoke",
		Usage:       "Invoke an entry point of server code",
		Description: "args: <entry-name> [version]",
		Action: func(c *cli.Context) {
			if len(c.Args()) > 2 || len(c.Args()) == 0 {
				cli.ShowCommandHelp(c, c.Command.Name)
				os.Exit(ExitIllegalNumberOfArgs)
			}
			version := "current"
			if len(c.Args()) == 2 {
				version = c.Args()[1]
			}
			InvokeServerCode(c.Args()[0], version)
		},
	},
	{
		Name:  "servercode:activate",
		Usage: "Activate a version",
		Action: func(c *cli.Context) {
			ShowCommandHelp(1, c)
			ActivateServerCode(c.Args()[0])
		},
	},
	{
		Name:  "servercode:delete",
		Usage: "Delete a version of server code",
		Action: func(c *cli.Context) {
			ShowCommandHelp(1, c)
			DeleteServerCode(c.Args()[0])
		},
	},
	{
		Name:        "servercode:hook-deploy",
		Usage:       "Delopy a hook config",
		Description: "args: <hooo-config-path> <version>",
		Action: func(c *cli.Context) {
			ShowCommandHelp(2, c)
			DeployHookConfig(c.Args()[0], c.Args()[1])
		},
	},
}
