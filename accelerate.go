package main

import (
	"crypto/tls"
	"errors"
	"log"
	"net"

	"github.com/fatih/color"
	"golang.org/x/net/http2"
)

var accTransp *http2.Transport

func checkErr(err error, str string) {
	if err != nil {
		log.Fatalln("Error: ", err, "context:", str)
	} else {
		log.Println(color.WhiteString("No error:"), str)
	}
}

func setupTransport() {
	accTransp = &http2.Transport{
		AllowHTTP: true,
		// Pretend we are dialing a TLS endpoint.
		// Note, we ignore the passed tls.Config
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return nil, errors.New("Should not be called")
		},
	}
}
