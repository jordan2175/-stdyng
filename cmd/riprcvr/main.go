package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

var proto = flag.Int("proto", 134, "Protocol number in IPv4 header")

func main() {
	flag.Parse()

	c, err := net.ListenPacket(fmt.Sprintf("ip4:%d", *proto), "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	p, err := ipv4.NewRawConn(c)
	if err != nil {
		log.Fatal(err)
	}
	p.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true)

	b := make([]byte, 8192)
	for {
		h, pp, cm, err := p.ReadFrom(b)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("%d bytes rcvd, %v, %v\n", len(pp), h, cm)
	}
}
