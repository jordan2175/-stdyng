package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	f, err := ioutil.TempFile("", "aflocal")
	if err != nil {
		log.Fatal(err)
	}
	addr := f.Name()
	os.Remove(f.Name())

	ln, err := net.Listen("unix", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(ln.Addr().String())
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				break
			}
			log.Println("local:", c.LocalAddr(), "remote:", c.RemoteAddr())
			c.Close()
		}
	}()

	(&http.Client{
		Transport: &http.Transport{
			Dial: func(network, address string) (net.Conn, error) {
				return (&net.Dialer{}).Dial("unix", ln.Addr().String())
			},
		},
	}).Get("http://localhost" + ln.Addr().String())
}
