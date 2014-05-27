package main

import (
	"log"
	"net"
)

func main() {
	for _, host := range []string{"www.kame.net", "www.google.com", "golang.org"} {
		conventionalWay(host)
		happyEyeballs(host)
	}
}

// See RFC 3493.
func conventionalWay(host string) {
	addrs, err := net.LookupHost(host) // almost equivalent to getaddrinfo
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range addrs {
		c, err := net.Dial("tcp", net.JoinHostPort(a, "80")) // almost equivalent to socket+bind+connect
		if err != nil {
			log.Println("oops", err)
			continue
		}
		log.Println("connected to", c.RemoteAddr(), "via", c.LocalAddr())
		c.Close()
	}
}

// See RFC 6555.
func happyEyeballs(host string) {
	d := net.Dialer{DualStack: true}
	c, err := d.Dial("tcp", net.JoinHostPort(host, "80"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to", c.RemoteAddr(), "via", c.LocalAddr())
	c.Close()
}
