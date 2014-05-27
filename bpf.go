package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"unsafe"
)

var ipv6OverEthernet = []syscall.BpfInsn{
	// make sure this is an IPv6 packet.
	*syscall.BpfStmt(syscall.BPF_LD+syscall.BPF_H+syscall.BPF_ABS, 12),
	*syscall.BpfJump(syscall.BPF_JMP+syscall.BPF_JEQ+syscall.BPF_K, 0x86dd, 0, 1),
	// if we passed all the tests, ask for the whole packet.
	*syscall.BpfStmt(syscall.BPF_RET+syscall.BPF_K, -1),
	// otherwise, drop it.
	*syscall.BpfStmt(syscall.BPF_RET+syscall.BPF_K, 0),
}

func bpfFile() (*os.File, error) {
	for i := 0; i < 10; i++ {
		f, err := os.OpenFile(fmt.Sprintf("/dev/bpf%d", i), os.O_RDWR, 0666)
		if err == nil {
			return f, nil
		}
	}
	return nil, syscall.ENOENT
}

func bpfInterface() (*net.Interface, error) {
	ift, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ifi *net.Interface
	avail := net.FlagUp | net.FlagBroadcast | net.FlagMulticast
	for _, i := range ift {
		if i.Flags&avail == avail {
			ifat, err := i.Addrs()
			if err != nil {
				return nil, err
			}
			if len(ifat) != 0 {
				ifi = &i
				break
			}
		}
	}
	if ifi == nil {
		return nil, syscall.ENOENT
	}
	return ifi, nil
}

func prepareBPF(fd int, name string) (int, error) {
	if err := syscall.SetBpfInterface(fd, name); err != nil {
		return 0, err
	}
	if err := syscall.CheckBpfVersion(fd); err != nil {
		return 0, err
	}
	if err := syscall.SetBpfImmediate(fd, 1); err != nil {
		return 0, err
	}
	if err := syscall.SetBpfPromisc(fd, 1); err != nil {
		return 0, err
	}
	buflen, err := syscall.BpfBuflen(fd)
	if err != nil {
		return 0, err
	}
	if _, err = syscall.BpfHeadercmpl(fd); err != nil {
		return 0, err
	}
	if err := syscall.SetBpfHeadercmpl(fd, 0); err != nil {
		return 0, err
	}
	if _, err := syscall.BpfTimeout(fd); err != nil {
		return 0, err
	}
	tv := syscall.Timeval{Usec: 10}
	if err := syscall.SetBpfTimeout(fd, &tv); err != nil {
		return 0, err
	}
	if err := syscall.SetBpf(fd, ipv6OverEthernet); err != nil {
		return 0, nil
	}
	if err := syscall.FlushBpf(fd); err != nil {
		return 0, err
	}
	return buflen, nil
}

func main() {
	if os.Getuid() != 0 {
		log.Fatal("must be run with administrator privileges")
	}

	f, err := bpfFile()
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	ifi, err := bpfInterface()
	if err != nil {
		log.Fatal(err)
	}
	blen, err := prepareBPF(int(f.Fd()), ifi.Name)
	if err != nil {
		log.Fatal(err)
	}
	var bpfh *syscall.BpfHdr
	b := make([]byte, blen)
	for {
		n, err := syscall.Read(int(f.Fd()), b)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			continue
		}
		bpfh = (*syscall.BpfHdr)(unsafe.Pointer(&b[0]))
		log.Printf("%+v\n", bpfh)
	}
}
