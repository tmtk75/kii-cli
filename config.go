package kiicli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/tmtk75/cli"
	"github.com/vaughan0/go-ini"
)

type GlobalConfig struct {
	AppId        string
	AppKey       string
	ClientId     string
	ClientSecret string
	Token        string
	Site         string
	endpointUrl  string
	devlogUrl    string
	Curl         bool
	SuppressExit bool
	UTC          bool
	usePName     bool
	profileName  string
}

const (
	ExitGeneralReason       = 1
	ExitIllegalNumberOfArgs = 2
	ExitNotLoggedIn         = 3
	ExitMissingParams       = 4
)

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
	p := Profile()
	host := hosts[p.Site]
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
	p := Profile()
	host := hosts[p.Site]
	if host == "" {
		print("missing site, use --site or set KII_SITE\n")
		os.Exit(ExitMissingParams)
	}
	return fmt.Sprintf("wss://%s:443/logs", host)
}

func (self *GlobalConfig) HttpHeaders(contentType string) map[string]string {
	p := Profile()
	m := map[string]string{
		"x-kii-appid":  p.AppId,
		"x-kii-appkey": p.AppKey,
	}
	if len(contentType) > 0 {
		m["content-type"] = contentType
	}
	return m
}

func (self *GlobalConfig) HttpHeadersWithAuthorization(contentType string) map[string]string {
	m := self.HttpHeaders(contentType)
	oauth2 := (&OAuth2Response{}).Load()

	// overwirte token with given value
	token := oauth2.AccessToken
	if self.Token != "" {
		token = self.Token
	}

	m["authorization"] = fmt.Sprintf("Bearer %s", token)
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
	d := DirPath([]string{dir})
	return d.MetaFilePath(filename)
}

type DirPath []string

func (dir DirPath) MetaFilePath(filename string) string {
	// Fix dir to take care of multi site & same app-id
	if globalConfig != nil && globalConfig.usePName && globalConfig.profileName != "" {
		dir = DirPath(append([]string{globalConfig.profileName}, dir...))
	}

	homedir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("%v", err)
	}
	confdirpath := path.Join(homedir, ".kii", path.Join(dir...))
	err = os.MkdirAll(confdirpath, os.ModeDir|0700)
	if err != nil {
		log.Fatalf("%v", err)
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

func loadIniFile(configPath string) (*ini.File, bool /* true: generated config */) {
	if b, _ := exists(configPath); !b {
		ioutil.WriteFile(configPath, []byte(_config), 0600)
		return nil, true
	}
	file, err := ini.LoadFile(configPath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &file, false
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

func profilePath(c *cli.Context) string {
	p := metaFilePath(".", "config")
	if c.String("profile-path") != "" {
		return c.String("profile-path")
	}
	return p
}

func SetupFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "app-id", Value: "", Usage: "AppID"},
		cli.StringFlag{Name: "app-key", Value: "", Usage: "AppKey"},
		cli.StringFlag{Name: "client-id", Value: "", Usage: "ClientID"},
		cli.StringFlag{Name: "client-secret", Value: "", Usage: "ClientSecret"},
		cli.StringFlag{Name: "token", Value: "", Usage: "Token to be used"},
		cli.StringFlag{Name: "site", Value: "", Usage: "us,jp,cn,sg"},
		cli.StringFlag{Name: "endpoint-url", Value: "", Usage: "Site URL"},
		cli.StringFlag{Name: "log-url", Value: "", Usage: "Log URL"},
		cli.BoolFlag{Name: "verbose", Usage: "Verbosely"},
		cli.StringFlag{Name: "profile,p", Value: DEFAULT_PROFILE, Usage: "Profile name for ~/.kii/config"},
		cli.StringFlag{Name: "profile-path", Usage: "Profile path instead of ~/.kii/config"},
		cli.BoolFlag{Name: "curl", Usage: "Print curl command saving body as a tmp file if body exists"},
		cli.BoolFlag{Name: "suppress-exit", Usage: "Suppress exit with 1 when receiving status code other than 2xx"},
		cli.StringFlag{Name: "http-proxy", Usage: "HTTP proxy URL to be used"},
		cli.BoolFlag{Name: "disable-http-proxy", Usage: "Disable HTTP proxy in your profile"},
		cli.BoolFlag{Name: "use-utc", Usage: "Format time in UTC"},
		cli.BoolFlag{Name: "use-profile-name", Usage: "Use profile name as config dirname"},
	}

	app.Before = func(c *cli.Context) error {
		// Setup logger
		if c.Bool("verbose") {
			logger = log.New(os.Stderr, "", log.LstdFlags)
		}

		profilePath := profilePath(c)
		logger.Printf("profile-path: %v", profilePath)

		inifile, gen := loadIniFile(profilePath)
		if gen {
			print(fmt.Sprintf("%v was created. please fill it with your credentials.\n", profilePath))
			os.Exit(ExitGeneralReason)
		}

		profile, _ := inifile.Get("", "profile")
		if profile == "" {
			profile = DEFAULT_PROFILE
		}
		if optProf := c.GlobalString("profile"); optProf != DEFAULT_PROFILE {
			profile = optProf
		}
		logger.Printf("profile: %v", profile)

		if profile != DEFAULT_PROFILE && len((*inifile)[profile]) == 0 {
			print(fmt.Sprintf("profile %s is not found in %v\n", profile, profilePath))
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
			Token:        getConf("token", "KII_TOKEN", ""),
			Site:         getConf("site", "KII_SITE", "site"),
			endpointUrl:  getConf("endpoint-url", "KII_ENDPOINT_URL", "endpoint_url"),
			devlogUrl:    getConf("log-url", "KII_LOG_URL", "log_url"),
			Curl:         c.GlobalBool("curl"),
			SuppressExit: c.GlobalBool("suppress-exit"),
			UTC:          c.GlobalBool("use-utc"),
			usePName:     c.GlobalBool("use-profile-name"),
			profileName:  profile,
		}

		proxy := c.String("http-proxy")
		if proxy == "" && !c.Bool("disable-http-proxy") {
			p, _ := inifile.Get(profile, "http_proxy")
			proxy = p
			if proxy == "" {
				p, _ := inifile.Get("", "http_proxy")
				proxy = p
			}
		}
		if proxy != "" {
			logger.Printf("http_proxy: %v", proxy)
			os.Setenv("HTTP_PROXY", proxy)
		}

		logger.Printf("dirname-to-store: %v\n", metaFilePath(globalConfig.AppId, "."))

		return nil
	}
}

func Flatten(a []cli.Command) []cli.Command {
	b := make([]cli.Command, 0, 16)
	for _, v := range a {
		if v.Subcommands == nil {
			b = append(b, v)
		} else {
			for _, i := range v.Subcommands {
				i.Name = v.Name + ":" + i.Name
				b = append(b, i)
			}
		}
	}
	return b
}
