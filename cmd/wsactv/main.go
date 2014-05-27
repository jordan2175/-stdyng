package main

import (
	"flag"
	"log"

	"golang.org/x/net/websocket"
)

var addr = flag.String("addr", "[fe80::1%25lo0]:5963", "target address")

func main() {
	flag.Parse()
	c, err := websocket.Dial("ws://"+*addr+"/reflect", "", "http://localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	if _, err := c.Write([]byte("HELLO-R-U-THERE")); err != nil {
		log.Fatal(err)
	}
	b := make([]byte, 1500)
	n, err := c.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(b[:n]))
}
