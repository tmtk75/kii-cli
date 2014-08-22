package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/vaughan0/go-ini"
)

type GlobalConfig struct {
	AppId        string
	AppKey       string
	ClientId     string
	ClientSecret string
	Site         string
	endpointUrl  string
	devlogUrl    string
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
	return fmt.Sprintf("https://%s/api", host)
}

func (self *GlobalConfig) EndpointUrlForApiLog() string {
	if self.devlogUrl != "" {
		return self.devlogUrl
	}
	hosts := map[string]string{
		"us": "apilog.kii.com",
		"jp": "apilog-jp.kii.com",
		"cn": "apilog-cn2.kii.com",
		"sg": "apilog-sg.kii.com",
	}
	host := hosts[globalConfig.Site]
	if host == "" {
		print("missing site, use --site or set KII_SITE\n")
		os.Exit(ExitMissingParams)
	}
	return fmt.Sprintf("wss://%s:443/logs", host)
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
func metaFilePath(dir string, filename string) string {
	homedir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	confdirpath := path.Join(homedir, ".kii", dir)
	err = os.MkdirAll(confdirpath, os.ModeDir|0700)
	if err != nil {
		panic(err)
	}
	return path.Join(confdirpath, filename)
}

var globalConfig *GlobalConfig
var _config = `# You can configure default profile, or --profile option is available
# profile = jp

[default]
app_id =
app_key =
client_id =
client_secret =
site = us

[jp]
app_id =
app_key =
client_id =
client_secret =
site = jp
`

func loadIniFile() *ini.File {
	configPath := metaFilePath(".", "config")
	if b, _ := exists(configPath); !b {
		ioutil.WriteFile(configPath, []byte(_config), 0600)
	}
	file, err := ini.LoadFile(configPath)
	if err != nil {
		panic(err)
	}
	return &file
}

func pickup(a ...string) string {
	for _, s := range a {
		if s != "" {
			return s
		}
	}
	return ""
}

const DEFAULT_PROFILE = "default"

func setupFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "app-id", Value: "", Usage: "AppID"},
		cli.StringFlag{Name: "app-key", Value: "", Usage: "AppKey"},
		cli.StringFlag{Name: "client-id", Value: "", Usage: "ClientID"},
		cli.StringFlag{Name: "client-secret", Value: "", Usage: "ClientSecret"},
		cli.StringFlag{Name: "site", Value: "", Usage: "us,jp,cn,sg"},
		cli.StringFlag{Name: "endpoint-url", Value: "", Usage: "Site URL"},
		cli.BoolFlag{Name: "verbose", Usage: "Verbosely"},
		cli.StringFlag{Name: "profile", Value: DEFAULT_PROFILE, Usage: "Profile name for ~/.kii/config"},
	}

	app.Before = func(c *cli.Context) error {
		// Setup logger
		if c.Bool("verbose") {
			logger = log.New(os.Stderr, "", log.LstdFlags)
		}

		inifile := loadIniFile()
		profile, _ := inifile.Get("", "profile")
		if optProf := c.GlobalString("profile"); optProf != DEFAULT_PROFILE {
			profile = optProf
		}
		logger.Printf("profile: %v", profile)

		if profile != DEFAULT_PROFILE && len((*inifile)[profile]) == 0 {
			print(fmt.Sprintf("profile %s is not found in ~/.kii/config\n", profile))
			os.Exit(ExitMissingParams)
		}

		getConf := func(gn, en, un string) string {
			ev := os.ExpandEnv("${" + en + "}")
			uv, _ := inifile.Get(profile, un)
			return pickup(c.GlobalString(gn), ev, uv)
		}

		globalConfig = &GlobalConfig{
			AppId:        getConf("app-id", "KII_APP_ID", "app_id"),
			AppKey:       getConf("app-key", "KII_APP_KEY", "app_key"),
			ClientId:     getConf("client-id", "KII_CLIENT_ID", "client_id"),
			ClientSecret: getConf("client-secret", "KII_CLIENT_SECRET", "client_secret"),
			Site:         getConf("site", "KII_SITE", "site"),
			endpointUrl:  getConf("endpoint-url", "KII_ENDPOINT_URL", "endpoint_url"),
			devlogUrl:    getConf("log-url", "KII_LOG_URL", "log_url"),
		}

		return nil
	}
}
