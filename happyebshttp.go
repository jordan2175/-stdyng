package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
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
		tr := &http.Transport{TLSClientConfig: &tls.Config{ServerName: host}}
		client := http.Client{Transport: tr}
		resp, err := client.Get(fmt.Sprintf("https://%s/", net.JoinHostPort(a, "443")))
		if err != nil {
			log.Println("oops", err)
			continue
		}
		log.Println("connected to", resp.Request.URL.Host)
		resp.Body.Close()
	}
}

// See RFC 6555.
func happyEyeballs(host string) {
	tr := &http.Transport{
		Dial:            (&net.Dialer{DualStack: true}).Dial,
		TLSClientConfig: &tls.Config{ServerName: host},
	}
	client := http.Client{Transport: tr}
	resp, err := client.Get(fmt.Sprintf("https://%s/", net.JoinHostPort(host, "443")))
	if err != nil {
		log.Println("oops", err)
		return
	}
	log.Println("connected to", resp.Request.URL.Host)
	resp.Body.Close()
}
