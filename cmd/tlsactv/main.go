// +build ignore

package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	tr := &http.Transport{
		Dial:            (&net.Dialer{Timeout: 3 * time.Second}).Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr}
	resp, err := client.Get("https://www.google.com/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	for _, cert := range resp.TLS.PeerCertificates {
		log.Println("version", cert.Version)
	}
}
