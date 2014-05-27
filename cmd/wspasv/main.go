package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var addr = flag.String("addr", ":5963", "server address")

func main() {
	flag.Parse()
	http.Handle("/reflect", websocket.Handler(reflector))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

func reflector(c *websocket.Conn) {
	defer c.Close()
	io.Copy(c, c)
}
