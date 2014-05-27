package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

var addr = flag.String("addr", "localhost", "target address")

func main() {
	flag.Parse()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         *addr,
			NextProtos:         []string{"foo", "bar", "baz"},
		},
	}
	client := http.Client{Transport: tr}
	resp, err := client.Get(fmt.Sprintf("https://%s/", net.JoinHostPort(*addr, "443")))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.Body)
	resp.Body.Close()
}
