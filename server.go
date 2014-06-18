package main

import (
	"code.google.com/p/go.net/websocket"
	"io"
	"log"
	"net/http"
)

func echoHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
