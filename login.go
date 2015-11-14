package kiicli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/tmtk75/cli"
)

const (
	OPENIDCONNECT        = "openidconnect.simple.discovery_document_url"
	FEDAUTH_SIGNUP_URI   = "federated-auth.signup-uri"
	FEDAUTH_SITE_URI     = "federated-auth.site-uri"
	FEDAUTH_REDIRECT_URI = "federated-auth.redirect-uri"
)

type OAuth2Request struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (self *GlobalConfig) OAuth2Request() *OAuth2Request {
	return &OAuth2Request{
		self.ClientId,
		self.ClientSecret,
	}
}

type OAuth2Response struct {
	Id          string `json:"id"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (self *OAuth2Response) Bytes() []byte {
	it, err := json.Marshal(self)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return it
}

func (self *OAuth2Response) Save(filename string) {
	err := ioutil.WriteFile(filename, self.Bytes(), 0600)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func (self *OAuth2Response) Load() *OAuth2Response {
	path := adminTokenFilePath()
	b, _ := exists(path)
	if !b {
		print("You've not logged in, first `auth login`\n")
		os.Exit(ExitNotLoggedIn)
	}
	self.LoadFrom(path)
	return self
}

func (self *OAuth2Response) LoadFrom(path string) *OAuth2Response {
	logger.Printf("Load %s\n", path)

	file, _ := os.Open(path)
	body, err := ioutil.ReadAll(bufio.NewReader(file))
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()
	json.Unmarshal(body, self)
	return self
}

func (self *OAuth2Response) Decode(res *HttpResponse) *OAuth2Response {
	d := json.NewDecoder(res.Body)
	err := d.Decode(&self)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return self
}

func retrieveAppAdminAccessToken() *OAuth2Response {
	p := Profile()
	token := p.OAuth2Request()
	headers := p.HttpHeaders("application/json")
	res := HttpPostJson("/oauth2/token", headers, token)
	oauth2res := &OAuth2Response{}
	return oauth2res.Decode(res)
}

func adminTokenFilePath() string {
	p := Profile()
	return metaFilePath(path.Join(".", p.AppId), "token")
}

func LoginAsAppAdmin(force bool) {
	if b, _ := exists(adminTokenFilePath()); b && !force {
		print("Already logged in, use `--force` to login\n")
		os.Exit(0)
	}
	res := retrieveAppAdminAccessToken()
	res.Save(adminTokenFilePath())
}

func IsMaster() bool {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/configuration/parameters/isMasterApp", p.AppId)
	headers := p.HttpHeadersWithAuthorization("text/plain")
	body := HttpGet(path, headers).Bytes()
	a, err := strconv.ParseBool(string(body))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return a
}

func ConfigureAsMaster() {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/configuration/parameters/isMasterApp", p.AppId)
	headers := p.HttpHeadersWithAuthorization("text/plain")
	r := bytes.NewBuffer([]byte("true"))
	body := HttpPut(path, headers, r).Bytes()
	fmt.Println(string(body))
}

func StepDownMaster() {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/configuration/parameters/isMasterApp", p.AppId)
	headers := p.HttpHeadersWithAuthorization("text/plain")
	r := bytes.NewBuffer([]byte("false"))
	body := HttpPut(path, headers, r).Bytes()
	fmt.Println(string(body))
}

// https://wiki.kii.com/display/Products/Federated+authentication
func ProvisionSlaveApp(appId string) (id, secret string) {
	p := Profile()
	if !IsMaster() {
		log.Fatalf("%v is not a master app", p.AppId)
	}
	path := fmt.Sprintf("/apps/%v/oauth2/clients", p.AppId)
	headers := p.HttpHeadersWithAuthorization("application/vnd.kii.Oauth2ClientCreationRequest+json")
	reduri := RedirectURI(appId)
	logger.Printf("externalID: %v\n", appId)
	logger.Printf("redirectURI: %v\n", reduri)

	type T struct {
		ExternalID  string `json:"externalID"`
		RedirectURI string `json:"redirectURI"`
	}
	t := T{ExternalID: appId, RedirectURI: reduri}
	j, _ := json.Marshal(t)
	body := HttpPost(path, headers, bytes.NewReader(j)).Bytes()

	type C struct {
		ClientID     string `json:"clientID"`
		ClientSecret string `json:"clientSecret"`
	}
	var c C
	err := json.Unmarshal(body, &c)
	if err != nil {
		log.Fatal("%v", err)
	}
	return c.ClientID, c.ClientSecret
}

func ConfigureSlaveApp(key, secret, masterAppId string) {
	p := Profile()
	if IsMaster() {
		log.Fatalf("%v is a master app", p.AppId)
	}
	path := fmt.Sprintf("/apps/%v", p.AppId)
	headers := p.HttpHeadersWithAuthorization("application/vnd.kii.AppModificationRequest+json")
	var t struct {
		SocialAuth struct {
			ConsumerKey    string `json:"kii.consumer_key"`
			ConsumerSecret string `json:"kii.consumer_secret"`
			MasterAppID    string `json:"kii.master_app_id"`
			MasterAppSite  string `json:"kii.master_app_site"`
		} `json:"socialAuth"`
	}
	t.SocialAuth.ConsumerKey = key
	t.SocialAuth.ConsumerSecret = secret
	t.SocialAuth.MasterAppID = masterAppId
	t.SocialAuth.MasterAppSite = FederatedAuthSiteURI(masterAppId)
	//
	logger.Printf("kii.consumer_key: %s\n", t.SocialAuth.ConsumerKey)
	logger.Printf("kii.consumer_secret: %s\n", t.SocialAuth.ConsumerSecret)
	logger.Printf("kii.master_app_id: %s\n", t.SocialAuth.MasterAppID)
	logger.Printf("kii.master_app_site: %s\n", t.SocialAuth.MasterAppSite)

	j, _ := json.Marshal(t)

	//logger.Printf("%v", string(j))
	body := HttpPost(path, headers, bytes.NewReader(j)).Bytes()
	fmt.Printf("%v\n", string(body))
}
func ConfigureOpenIDConnect(key, secret, slaveAppID string) {
	p := Profile()
	if !IsMaster() {
		log.Fatalf("%v is not a master app", p.AppId)
	}
	path := fmt.Sprintf("/apps/%v/configuration/parameters", p.AppId)
	headers := p.HttpHeadersWithAuthorization("application/vnd.kii.appconfigparamsmodificationrequest+json")
	var t struct {
		ConsumerKey          string `json:"openidconnect.simple.consumer_key"`
		ConsumerSecret       string `json:"openidconnect.simple.consumer_secret"`
		DiscoveryDocumentURL string `json:"openidconnect.simple.discovery_document_url"`
		Scope                string `json:"openidconnect.simple.scope"`
	}
	t.ConsumerKey = key
	t.ConsumerSecret = secret
	t.DiscoveryDocumentURL = DiscoveryDocumentURL(slaveAppID)
	t.Scope = "openid email"
	//
	logger.Printf("openidconnect.simple.consumer_key: %s\n", t.ConsumerKey)
	logger.Printf("openidconnect.simple.consumer_secret: %s\n", t.ConsumerSecret)
	logger.Printf("openidconnect.simple.discovery_document_url: %s\n", t.DiscoveryDocumentURL)
	logger.Printf("openidconnect.simple.scope: %s\n", t.Scope)

	j, _ := json.Marshal(t)

	//logger.Printf("%v", string(j))
	body := HttpPatch(path, headers, bytes.NewReader(j)).Bytes()
	fmt.Printf("%v\n", string(body))
}

func BuildSignUpURL() string {
	p := Profile()
	if IsMaster() {
		log.Fatalf("%v is a master app", p.AppId)
	}
	return FederatedAuthSignUpURI(p.AppId)
}

func ShowSlaveInfo(clientID string) {
	p := Profile()
	path := fmt.Sprintf("/apps/%s/oauth2/clients/%s", p.AppId, clientID)
	headers := p.HttpHeadersWithAuthorization("application/json")
	body := HttpGet(path, headers).Bytes()
	fmt.Printf("%v\n", string(body))
}

func GenerateCert() {
	p := Profile()
	if !IsMaster() {
		log.Fatalf("%v is not a master app", p.AppId)
	}
	path := fmt.Sprintf("/apps/%s/oauth2/certs", p.AppId)
	headers := p.HttpHeadersWithAuthorization("application/json")
	r := strings.NewReader("{}")
	body := HttpPost(path, headers, r).Bytes()
	fmt.Printf("%v\n", string(body))
}

func ConfigureFederatedAuth(pname string) {
	p := Profile()

	//
	if p.profileName == pname {
		log.Fatalf("You must specify different profile name for %v\n", pname)
	}

	// beehive-dev-parent auth configure-as-master
	ConfigureAsMaster()
	fmt.Printf("%v became master\n", p.AppId)

	// beehive-dev-parent auth generate-cert
	GenerateCert()
	fmt.Printf("Generated a cert\n")

	// beehive-dev-parent auth provision-slave 20479a45
	slave := FindAppID(pname)
	id, secret := ProvisionSlaveApp(slave)
	fmt.Printf("Provisioned the master with a slave app, %v\n", slave)

	args := []string{}
	if p.Verbose {
		args = append(args, "--verbose")
	}
	if p.Curl {
		args = append(args, "--curl")
	}

	// beehive-dev-child auth configure-as-slave a9f0f99uibd0jeja6rn8o4mkjg899sjrph9m chn8rrm4hi7n6dal7pt79l1a06ukhvukk50a3bt1dlnscul7ohomg09ikb5493v 5bae5f53
	err := run("kii-cli", append(args, []string{"--profile", pname, "auth", "federated", "configure-as-slave", id, secret, p.AppId}...))
	if err != nil {
		log.Fatalf("%v", err)
	}

	// beehive-dev-child auth show-signup-url
	err = run("kii-cli", append(args, []string{"--profile", pname, "auth", "federated", "show-signup-url"}...))
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func run(s string, args []string) error {
	wd, _ := os.Getwd()
	logger.Printf("[debug] %v %v in %v\n", s, args, wd)
	//cmd := exec.Command(path.Join(wd, s), args...)
	cmd := exec.Command(s, args...)
	cmd.Dir = wd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return err
}

func kiiapps(siteCode, path string) (string, error) {
	u := map[string]string{
		"us":  "https://%s.us.kiiapps.com",
		"jp":  "https://%s.jp.kiiapps.com",
		"sg":  "https://%s.sg.kiiapps.com",
		"cn2": "https://%s.cn2.kiiapps.com",
		"cn3": "https://%s.cn3.kiiapps.com",
	}
	logger.Printf("site-code: %v", siteCode)
	if _, has := u[siteCode]; !has {
		return "", fmt.Errorf("Unknown site code: %v", siteCode)
	}
	return u[siteCode] + path, nil
}

func kiiapi(siteCode, path string) (string, error) {
	u := map[string]string{
		"us":  "https://api.kii.com/api",
		"jp":  "https://api-jp.kii.com",
		"sg":  "https://api-sg.kii.com",
		"cn2": "https://api-cn2.kii.com",
		"cn3": "https://api-cn3.kii.com",
	}
	logger.Printf("site-code: %v", siteCode)
	if _, has := u[siteCode]; !has {
		return "", fmt.Errorf("Unknown site code: %v", siteCode)
	}
	return u[siteCode] + path, nil
}

func definedRedirectURI(siteCode string) (string, error) {
	return kiiapps(siteCode, "/api/apps/%s/integration/webauth/callback")
}

func definedSiteURL(siteCode string) (string, error) {
	return kiiapps(siteCode, "/api")
}

func definedSignUpURL(siteCode string) (string, error) {
	return kiiapps(siteCode, "/api/apps/%s/integration/webauth/connect?id=kii")
}

func definedDiscoveryDocumentURL(siteCode string) (string, error) {
	return kiiapi(siteCode, "/apps/%s/.well-known/openid-configuration")
}

func RedirectURI(appId string) string {
	return findValueOf(FEDAUTH_REDIRECT_URI, appId, definedRedirectURI, 2)
}

func FederatedAuthSiteURI(appId string) string {
	return findValueOf(FEDAUTH_SITE_URI, appId, definedSiteURL, 1)
}

func DiscoveryDocumentURL(appId string) string {
	return findValueOf(OPENIDCONNECT, appId, definedDiscoveryDocumentURL, 1)
}

func FederatedAuthSignUpURI(appId string) string {
	return findValueOf(FEDAUTH_SIGNUP_URI, appId, definedSignUpURL, 2)
}

func findValueOf(key, appId string, defined func(string) (string, error), numargs int) string {
	args := make([]interface{}, numargs)
	for i := 0; i < numargs; i++ {
		args[i] = appId
	}
	s := FindIniFile(appId)
	if _, has := s[key]; !has {
		if _, has := s["site"]; !has {
			log.Fatalf("Found %v, but it doens't have %s. Please check your config file.\n", appId, key)
		}
		f, err := defined(s["site"])
		if err != nil {
			log.Fatalf("Unknown site code: %v of %v", s["site"], appId)
		}
		logger.Printf("Found %s for %v, %v", key, appId, f)
		return fmt.Sprintf(f, args...)

	}
	logger.Printf("Found %s for %v, %v", key, appId, s[key])
	u := s[key]
	return fmt.Sprintf(u, args...)
}

var LoginCommands = []cli.Command{
	{
		Name:  "login",
		Usage: "Login as AppAdmin",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "force", Usage: "Force to login"},
		},
		Action: func(c *cli.Context) {
			LoginAsAppAdmin(c.Bool("force"))
		},
	},
	{
		Name:  "info",
		Usage: "Print login info",
		Action: func(c *cli.Context) {
			res := &OAuth2Response{}
			res.Load()
			fmt.Println(res.AccessToken)
		},
	},
	{
		Name:        "federated",
		Usage:       "Federated auth",
		Subcommands: FederatedCommands,
	},
}

var FederatedCommands = []cli.Command{
	{
		Name:  "is-master",
		Usage: "Describe if mater",
		Action: func(c *cli.Context) {
			b := IsMaster()
			fmt.Println(b)
		},
	},
	{
		Name:  "configure-as-master",
		Usage: "Congfigure as master",
		Action: func(c *cli.Context) {
			ConfigureAsMaster()
		},
	},
	{
		Name:  "step-down-master",
		Usage: "Step down master",
		Action: func(c *cli.Context) {
			StepDownMaster()
		},
	},
	{
		Name:  "generate-cert",
		Usage: "Generate cert",
		Action: func(c *cli.Context) {
			GenerateCert()
		},
	},
	{
		Name:  "provision-slave",
		Usage: "Provision a master app for slave",
		Args:  "<slave-app-id>",
		Action: func(c *cli.Context) {
			s, _ := c.ArgFor("slave-app-id")
			id, secret := ProvisionSlaveApp(s)
			fmt.Printf("%s %s\n", id, secret)
		},
	},
	{
		Name:  "configure-as-slave",
		Usage: "Configure as slave",
		Args:  "<consumer-key> <consumer-secret> <master-app-id>",
		Action: func(c *cli.Context) {
			key, _ := c.ArgFor("consumer-key")
			secret, _ := c.ArgFor("consumer-secret")
			appId, _ := c.ArgFor("master-app-id")
			ConfigureSlaveApp(key, secret, appId)
		},
	},
	//{
	//	// https://wiki.kii.com/pages/viewpage.action?pageId=17006919
	//	Name:  "configure-openidconnect",
	//	Usage: "Configure OpenID connect for master",
	//	Args:  "<consumer-key> <consumer-secret> <slave-app-id>",
	//	Action: func(c *cli.Context) {
	//		key, _ := c.ArgFor("consumer-key")
	//		secret, _ := c.ArgFor("consumer-secret")
	//		appID, _ := c.ArgFor("slave-app-id")
	//		ConfigureOpenIDConnect(key, secret, appID)
	//	},
	//},
	//{
	//	Name:  "show-slave-info",
	//	Usage: "Print slave info",
	//	Args:  "<client-id>",
	//	Action: func(c *cli.Context) {
	//		clientID, _ := c.ArgFor("client-id")
	//		ShowSlaveInfo(clientID)
	//	},
	//},
	{
		Name:  "show-signup-url",
		Usage: "Print a URL to sign up with the configured master",
		Action: func(c *cli.Context) {
			s := BuildSignUpURL()
			fmt.Println(s)
		},
	},
	{
		Name:  "configure",
		Usage: "Configure master with a slave",
		Description: `
   FIXME
		`,
		Args: "<slave-profile-name>",
		Action: func(c *cli.Context) {
			pname, _ := c.ArgFor("slave-profile-name")
			ConfigureFederatedAuth(pname)
		},
	},
}
