package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/ipv4"
)

var (
	group = flag.String("group", "224.0.0.254", "ipv4 group address")
	port  = flag.Int("port", 12345, "udp service port")
)

func main() {
	flag.Parse()
	ip := net.ParseIP(*group)

	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		log.Fatal(err)
	}
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
	lsa := syscall.SockaddrInet4{Port: *port}
	copy(lsa.Addr[:], ip.To4())
	if err := syscall.Bind(s, &lsa); err != nil {
		syscall.Close(s)
		log.Fatal(err)
	}
	f := os.NewFile(uintptr(s), "")
	c, err := net.FilePacketConn(f)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	p := ipv4.NewPacketConn(c)
	defer p.Close()

	ift, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	avail := net.FlagMulticast | net.FlagUp
	for _, ifi := range ift {
		if ifi.Flags&avail != avail {
			continue
		}
		if err := p.JoinGroup(&ifi, &net.UDPAddr{IP: ip}); err != nil {
			log.Println(err, "on", ifi)
		}
	}

	log.Println(c.LocalAddr())
	go receiver(c)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-sig:
			os.Exit(0)
		}
	}
}

func receiver(c net.PacketConn) {
	b := make([]byte, 1500)
	for {
		n, peer, err := c.ReadFrom(b)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("%v bytes received from %v\n", n, peer)
	}
}
