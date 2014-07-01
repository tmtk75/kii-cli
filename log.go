package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
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
				fmt.Printf("%04d: %s\n", i, m.Log())
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

var LogCommands = []cli.Command{
	{
		Name:  "log",
		Usage: "Disply logs for an app",
		Action: func(c *cli.Context) {
			StartLogging()
		},
	},
}
