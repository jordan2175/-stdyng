package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

var (
	proto = flag.Int("proto", 134, "Protocol number in IPv4 header")
	dst   = flag.String("dst", "255.255.255.255", "destionation IPv4 address")
)

func main() {
	flag.Parse()
	ip := net.ParseIP(*dst)
	if ip == nil {
		log.Fatal("destination not found")
	}

	c, err := net.ListenPacket(fmt.Sprintf("ip4:%d", *proto), "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	p, err := ipv4.NewRawConn(c)
	if err != nil {
		log.Fatal(err)
	}

	b := []byte("HELLO-R-U-THERE")
	h := &ipv4.Header{
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TotalLen: ipv4.HeaderLen + len(b),
		Protocol: *proto,
		Dst:      ip.To4(),
	}
	if err := p.WriteTo(h, b, nil); err != nil {
		log.Println(err)
	}
}
