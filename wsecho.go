package main

import (
	"io"
	"log"
	"net/http"

	"github.com/tmtk75/cli"

	"code.google.com/p/go.net/websocket"
)

func echoHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func StartWSEchoServer() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Developing
var WSEchoCommands = []cli.Command{
	{
		Name:  "server",
		Usage: "WebSocket echo server for testing",
		Action: func(c *cli.Context) {
			StartWSEchoServer()
		},
	},
}
