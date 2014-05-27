package main

import (
	"crypto/tls"
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("This is a dumb ALPN server.\n"))
	})
	config := tls.Config{
		NextProtos:   []string{"foo", "bar", "baz"},
		Certificates: make([]tls.Certificate, 1),
		ServerName:   "localhost",
	}
	var err error
	config.Certificates[0], err = tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      ":443",
		TLSConfig: &config,
		Handler:   handler,
	}
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
