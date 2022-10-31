package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Alex-Eftimie/utils"
	"github.com/fatih/color"
	"github.com/soheilhy/cmux"
)

var proxyConn net.Conn
var s5s *Socks5Server
var h1s *HTTPServer

func main() {
	utils.DebugLevel = 1000
	var err error
	setupTransport()

	manageUpstream(Co.ProxyAddr)

	// if len(os.Args) > 1 {
	// os.Exit(0)
	// }
	addr := fmt.Sprintf("127.0.0.1:%d", Co.BindPort)
	cl, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to start: ", err)
	}
	mux := cmux.New(cl)

	// run the matchers
	socks5Matcher := mux.Match(Socks5Matcher())
	httpMatcher := mux.Match(cmux.HTTP1Fast())

	// initialize the servers
	s5s := newSocks5Server(addr)
	h1s := newHTTPServer(addr)
	// h1s.
	// run the servers
	go h1s.Serve(httpMatcher)
	go s5s.Serve(socks5Matcher)

	mux.Serve()
}

func manageUpstream(host string) {
	var err error
	var ctx = context.Background()

	for {
		proxyConn, err = net.Dial("tcp", host)
		if err == nil {
			break
		}
		log.Println(color.RedString("net.Dial error:"), err)
		log.Println("Retrying in 5 seconds")
		alert("Accelerator Error", fmt.Sprintf("Retrying in 5 seconds\n\nError: %s", err.Error()))
		time.Sleep(5 * time.Second)
	}
	h2cc, err = accTransp.NewClientConn(proxyConn)
	checkErr(err, "accTransp.NewClientConn")
	log.Println(color.GreenString("Connected to upstream:"), host)

	first := true
	go func() {
		for {
			time.Sleep(5 * time.Second)
			// keep alive
			err = h2cc.Ping(ctx)
			if err != nil {
				log.Println(color.RedString("Keepalive error:"), err, ", reconnecting")
				alert("Accelerator Error", fmt.Sprintf("Failed keepalive, reconnecting\n\nError: %s", err.Error()))
				go manageUpstream(host)
				return
			}
			log.Println(color.GreenString("Alive"))
			if first {

				notify("", fmt.Sprintf("Listening on port %d, forwarding to: %s", Co.BindPort, Co.ProxyAddr))
				first = false
			}
		}
	}()
}
