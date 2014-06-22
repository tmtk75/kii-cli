package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"code.google.com/p/go.crypto/ssh/terminal"
)

type Headers map[string]string

func DeployServerCode(serverCodePath string, activate bool) {
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
}

func InvokeServerCode(entryName string, version string) {
	path := fmt.Sprintf("/apps/%s/server-code/versions/%s/%s", globalConfig.AppId, version, entryName)
	headers := globalConfig.HttpHeadersWithAuthorization("application/json")
	var b []byte
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		b = HttpPost(path, headers, strings.NewReader("{}")).Bytes()
	} else {
		b = HttpPost(path, headers, os.Stdin).Bytes()
	}
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
