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
	"text/template"
	"time"

	"github.com/codegangsta/cli"

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

type Format map[string]string

var format Format

func LoadFormat() Format {
	path := "./format.json"
	e, err := exists(path)
	if err != nil {
		panic(err)
	}

	f := make(Format)
	if !e {
		logger.Printf("%v is missing", path)
		return f
	}

	body, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &f); err != nil {
		panic(err)
	}

	for k, v := range f {
		f[k] = convertLogFormat(v)
	}

	return f
}

func convertLogFormat(f string) string {
	re, _ := regexp.Compile("\\${[a-zA-Z-_]+}")
	k := re.ReplaceAllFunc([]byte(f), func(a []byte) []byte {
		return []byte(fmt.Sprintf("{{.%v}}", string(a[2:len(a)-1])))
	})
	return string(k)
}

func (m *RawLog) Print(idx int) {
	key := (*m)["key"].(string)
	f := format[key]
	if f == "" {
		return
	}

	t, _ := template.New(key).Parse(f)
	w := bytes.NewBuffer([]byte{})
	t.Execute(w, *m)
	fmt.Printf("%v\n", w)
}

var LogCommands = []cli.Command{
	{
		Name:  "log",
		Usage: "Disply logs for an app",
		Action: func(c *cli.Context) {
			format = LoadFormat()
			StartLogging()
		},
	},
}
