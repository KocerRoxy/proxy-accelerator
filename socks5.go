package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/Alex-Eftimie/netutils"
	"github.com/Alex-Eftimie/socks5"
	"github.com/Alex-Eftimie/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/soheilhy/cmux"
)

// Socks5Server handles socks5 proxy
type Socks5Server struct {
	*socks5.Server
}

// InternalHTTPWriter is used by socks5 to pipe content
type InternalHTTPWriter struct {
	net.Conn
}

// Hijack returns the inner net.Conn
func (m InternalHTTPWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return m.Conn, nil, nil
}

// Header returns a Internal http.Header
func (InternalHTTPWriter) Header() http.Header { return make(http.Header) }

// WriteHeader does nothing
func (InternalHTTPWriter) WriteHeader(int) {}

// Internal returns true if it's a Internal request
func (InternalHTTPWriter) Internal() bool { return true }

// Internal interface will be used to identify internal requests in http
type Internal interface {
	Internal() bool
}

// Socks5Matcher helps cmux determine if a request is socks5
func Socks5Matcher() cmux.Matcher {
	return func(r io.Reader) bool {
		b := make([]byte, 1)
		r.Read(b)
		return b[0] == 0x05
	}
}

func newSocks5Server(addr string) *Socks5Server {
	s5s := &Socks5Server{}
	s5s.Server = &socks5.Server{
		Addr:          addr,
		TunnelHandler: s5s.HandleTunnel,
	}
	return s5s
}

// HandleTunnel is called when a socks5 client requests to connect
func (s *Socks5Server) HandleTunnel(uinfo *netutils.UserInfo, ip string, c net.Conn, upstreamHost string, upstreamPort int, sc socks5.StatusCallback) {
	color.Yellow(spew.Sdump("[Socks5] Connect", upstreamHost, uinfo))

	w := &InternalHTTPWriter{c}

	r, err := http.NewRequest("CONNECT", fmt.Sprintf("http://%s:%d", upstreamHost, upstreamPort), utils.PrintReader{Reader: c, Prefix: "from HTTP"})
	r.RemoteAddr = "Socks5Server"
	if err != nil {
		log.Println("Error at http.NewRequest,", err)
		sc("general-failure.status", socks5.StatusGeneralFailure)
		return
	}
	sc("succeeded.status", socks5.StatusSucceeded)

	h1s.ServeHTTP(w, r)
}
