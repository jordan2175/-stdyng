// +build linux

package main

import (
	"log"
	"net"
	"syscall"
)

func main() {
	const proto = (syscall.ETH_P_ARP<<8)&0xff00 | syscall.ETH_P_ARP>>8 // need to take care of machine dependent stuff such as endianness when we use syscall directly
	s, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, proto)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(s)
	if err := syscall.SetNonblock(s, true); err != nil {
		log.Fatal(err)
	}
	ifi, err := net.InterfaceByName("eth0")
	if err != nil {
		log.Fatal(err)
	}
	lla := syscall.SockaddrLinklayer{Protocol: proto, Ifindex: ifi.Index}
	if err := syscall.Bind(s, &lla); err != nil {
		log.Fatal(err)
	}
	b := make([]byte, 6+6+4+2+8192+4) // w/ FCS, w/o preamble, dot1ad/dot1ah tags, SFD and IFG
	for {
		var (
			n    int
			peer syscall.Sockaddr
		)
		for {
			n, peer, err = syscall.Recvfrom(s, b, 0)
			if err != nil {
				n = 0
				if err == syscall.EAGAIN {
					continue
				}
				log.Fatal(err)
			}
			log.Printf("%v received from %+v\n", n, peer)
		}
	}
}
