package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

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

func StartLogging() {
	req := &AuthReq{
		os.ExpandEnv("${KII_APP_ID}"),
		os.ExpandEnv("${KII_APP_KEY}"),
		os.ExpandEnv("${KII_CLIENT_ID}"),
		os.ExpandEnv("${KII_CLIENT_SECRET}"),
		"tail",
	}
	j, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s", string(j))

	ws, err := websocket.Dial("ws://apilog.kii.com:80/logs", "", "http://localhost/")
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
