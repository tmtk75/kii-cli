package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"

	"code.google.com/p/go.net/websocket"
)

type RawLog map[string]interface{}

func (self *RawLog) Log() *Log {
	f := time.RFC3339Nano // "2006-01-02T15:04:05.999Z"
	t, err := time.Parse(f, (*self)["time"].(string))
	if err != nil {
		panic(err)
	}
	return &Log{
		(*self)["key"].(string),
		(*self)["level"].(string),
		(*self)["description"].(string),
		t,
	}
}

type Log struct {
	Key         string
	Level       string
	Description string
	Time        time.Time
}

type AuthRequest struct {
	AppID        string `json:"appID"`
	AppKey       string `json:"appKey"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	//	Token        string
	Command string `json:"command"` // 'tail' or 'cat'
	//	UserID       string
	//	Level        string
	//	DateFrom     string
	//	DateTo       string
}

func (self *GlobalConfig) AuthRequest() *AuthRequest {
	req := &AuthRequest{
		self.AppId,
		self.AppKey,
		self.ClientId,
		self.ClientSecret,
		"tail",
	}
	return req
}

func StartLogging() {
	req := globalConfig.AuthRequest()
	j, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s", string(j))

	url := globalConfig.EndpointUrlForApiLog()
	logger.Printf("log-url: %s", url)
	logger.Printf("%s", j)
	ws, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		panic(err)
	}
	_, err = ws.Write(j)
	if err != nil {
		panic(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	go func() {
		_ = <-sigc
		ws.Close()
		os.Exit(0)
	}()

	rch := make(chan []RawLog)

	go func() {
		for {
			//log.Println("will read")
			msg := <-rch
			for i, m := range msg {
				m.Print(i)
			}
		}
	}()

	for {
		select {
		default:
			var msg []RawLog
			err = websocket.JSON.Receive(ws, &msg)
			if err == io.EOF {
				os.Exit(0)
			}
			if err != nil {
				panic(err)
			}
			rch <- msg
			//log.Printf("wrote %d", len(msg))
		}
	}
}

type RawFormat map[string]string
type Format map[string]*template.Template

var format Format

func LoadFormat(path string) Format {
	e, err := exists(path)
	if err != nil {
		panic(err)
	}

	if !e {
		logger.Printf("%v is missing", path)
		return make(Format)
	}

	body, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	f := make(RawFormat)
	if err := json.Unmarshal(body, &f); err != nil {
		panic(err)
	}

	r := make(Format)
	for k, v := range f {
		s := convertLogFormat(v)
		t, err := template.New(k).Parse(s)
		if err != nil {
			panic(err)
		}
		r[k] = t
	}

	return r
}

// Kii official format in nodejs is ${foobar}, golang template in std pkg is {{.foobar}}
func convertLogFormat(f string) string {
	re, _ := regexp.Compile("\\${[a-zA-Z-_]+}")
	k := re.ReplaceAllFunc([]byte(f), func(a []byte) []byte {
		s := norm(string(a[2 : len(a)-1]))
		return []byte(fmt.Sprintf("{{.%v}}", s))
	})
	return string(k)
}

// Kii official format may contain hyphens in key, but golang template cannot handle easily.
// This func normalizes key to be handled in golang template.
func norm(k string) string {
	return strings.Replace(k, "-", "_", -1)
}

func (m *RawLog) Print(idx int) {
	key := (*m)["key"].(string)
	f := format[key]
	if f == nil {
		fmt.Printf("%v\n", *m.Log())
		return
	}

	i := make(RawLog)
	for k, v := range *m {
		i[norm(k)] = v
	}

	w := bytes.NewBuffer([]byte{})
	f.Execute(w, i)
	fmt.Printf("%v\n", w)
}

var LogCommands = []cli.Command{
	{
		Name:  "log",
		Usage: "Disply logs for an app",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "format-file", Usage: "File path to a format file", Value: (func() string {
				d, _ := homedir.Dir()
				return fmt.Sprintf("%v/.kii/format.json", d)
			})(),
			},
		},
		Action: func(c *cli.Context) {
			format = LoadFormat(c.String("format-file"))
			StartLogging()
		},
	},
}
