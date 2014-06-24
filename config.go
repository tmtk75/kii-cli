package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/codegangsta/cli"
	"github.com/vaughan0/go-ini"
)

type GlobalConfig struct {
	AppId        string
	AppKey       string
	ClientId     string
	ClientSecret string
	Site         string
	endpointUrl  string
}

func (self *GlobalConfig) EndpointUrl() string {
	if self.endpointUrl != "" {
		return self.endpointUrl
	}
	hosts := map[string]string{
		"us": "api.kii.com",
		"jp": "api-jp.kii.com",
		"cn": "api-cn2.kii.com",
		"sg": "api-sg.kii.com",
	}
	host := hosts[globalConfig.Site]
	if host == "" {
		print("missing site, use --site or set KII_SITE\n")
		os.Exit(ExitMissingParams)
	}
	return fmt.Sprintf("http://%s/api", host)
}

func (self *GlobalConfig) HttpHeaders(contentType string) map[string]string {
	m := map[string]string{
		"x-kii-appid":  globalConfig.AppId,
		"x-kii-appkey": globalConfig.AppKey,
	}
	if len(contentType) > 0 {
		m["content-type"] = contentType
	}
	return m
}

func (self *GlobalConfig) HttpHeadersWithAuthorization(contentType string) map[string]string {
	m := self.HttpHeaders(contentType)
	oauth2 := (&OAuth2Response{}).Load()
	m["authorization"] = fmt.Sprintf("Bearer %s", oauth2.AccessToken)
	return m
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Return ~/.kii/${filename}
func metaFilePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	confdirpath := path.Join(usr.HomeDir, ".kii")
	err = os.MkdirAll(confdirpath, os.ModeDir|0700)
	if err != nil {
		panic(err)
	}
	return path.Join(confdirpath, filename)
}

var globalConfig *GlobalConfig
var _config = `[default]
app_id =
app_key =
client_id =
client_secret =
site = us
`

func loadIniFile() *ini.File {
	configPath := metaFilePath("config")
	if b, _ := exists(configPath); !b {
		ioutil.WriteFile(configPath, []byte(_config), 0600)
	}
	file, err := ini.LoadFile(configPath)
	if err != nil {
		panic(err)
	}
	return &file
}

func pickup(a string, b string, c string) string {
	if len(a) > 0 {
		return a
	} else if len(b) > 0 {
		return b
	}
	return c
}

func setupFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{"app-id", "", "AppID"},
		cli.StringFlag{"app-key", "", "AppKey"},
		cli.StringFlag{"client-id", "", "ClientID"},
		cli.StringFlag{"client-secret", "", "ClientSecret"},
		cli.StringFlag{"site", "", "us,jp,cn,sg"},
		cli.StringFlag{"endpoint-url", "", "Site URL"},
		cli.BoolFlag{"verbose", "Verbosely"},
		cli.StringFlag{"profile", "default", "Profile name for ~/.kii/config"},
	}

	app.Before = func(c *cli.Context) error {
		profile := c.GlobalString("profile")
		inifile := loadIniFile()
		if profile != "default" && len((*inifile)[profile]) == 0 {
			print(fmt.Sprintf("profile %s is not found in ~/.kii/config\n", profile))
			os.Exit(ExitMissingParams)
		}

		get := func(name string) string {
			v, _ := inifile.Get(profile, name)
			return v
		}
		globalConfig = &GlobalConfig{
			pickup(c.GlobalString("app-id"), os.ExpandEnv("${KII_APP_ID}"), get("app_id")),
			pickup(c.GlobalString("app-key"), os.ExpandEnv("${KII_APP_KEY}"), get("app_key")),
			pickup(c.GlobalString("client-id"), os.ExpandEnv("${KII_CLIENT_ID}"), get("client_id")),
			pickup(c.GlobalString("client-secret"), os.ExpandEnv("${KII_CLIENT_SECRET}"), get("client_secret")),
			pickup(c.GlobalString("site"), os.ExpandEnv("${KII_SITE}"), get("site")),
			pickup(c.GlobalString("endpoint-url"), os.ExpandEnv("${KII_ENDPOINT_URL}"), get("endpoint_url")),
		}
		if c.Bool("verbose") {
			logger = &_Logger{}
		}
		return nil
	}
}
