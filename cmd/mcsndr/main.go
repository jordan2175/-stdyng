package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/ipv4"
)

var (
	group = flag.String("group", "224.0.0.254", "ipv4 group address")
	port  = flag.String("port", "12345", "udp service port")
)

func main() {
	flag.Parse()

	c, err := net.ListenPacket("udp", net.JoinHostPort(*group, "0"))
	if err != nil {
		log.Fatal(err)
	}
	p := ipv4.NewPacketConn(c)
	defer p.Close()
	dst, err := net.ResolveUDPAddr("udp", net.JoinHostPort(*group, *port))
	if err != nil {
		log.Fatal(err)
	}

	log.Println(c.LocalAddr())
	go sender(p, dst)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-sig:
			os.Exit(0)
		}
	}
}

func sender(p *ipv4.PacketConn, dst net.Addr) {
	ift, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; ; i++ {
		avail := net.FlagMulticast | net.FlagUp
		for _, ifi := range ift {
			if ifi.Flags&avail != avail {
				continue
			}
			p.SetMulticastInterface(&ifi)
			b := []byte(fmt.Sprintf("#%v HELLLO-R-U-THERE", i))
			if _, err := p.WriteTo(b, nil, dst); err != nil {
				log.Println(err, "on", ifi)
				continue
			}
			log.Printf("%v bytes sent to %v via %v\n", len(b), dst, ifi)
			time.Sleep(time.Second)
		}
	}
}
