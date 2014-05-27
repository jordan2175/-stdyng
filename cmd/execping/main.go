package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	environ := os.Environ()
	command := "/sbin/ping"
	args := []string{command, "-c", "5", "localhost"}

	cmd := exec.Command(command)
	cmd.Env = environ
	cmd.Args = args
	cmd.Dir = "."
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	b := make([]byte, 1024)
	for {
		n, err := stdout.Read(b)
		if err != nil {
			break
		}
		log.Print(string(b[:n]))
	}
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
