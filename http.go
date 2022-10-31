package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Alex-Eftimie/netutils"
	"github.com/fatih/color"
	"golang.org/x/net/http2"
)

var h2cc *http2.ClientConn

// HTTPServer handles http proxy
type HTTPServer struct {
	*http.Server
}

func newHTTPServer(addr string) *HTTPServer {
	hs := &HTTPServer{}
	hs.Server = &http.Server{
		Addr:    addr,
		Handler: hs,
	}
	return hs
}

func (cs *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[HTTP]", r.RemoteAddr, r.Method, r.URL)

	// r.Write(os.Stderr)

	var pipW *io.PipeWriter
	var pipR *io.PipeReader

	if r.Method == "CONNECT" {
		pipR, pipW = io.Pipe()
		r.Body = pipR
	}

	resp, err := h2cc.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, isInternal := w.(Internal)

	if r.Method == "CONNECT" {
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Tunneling not supported", http.StatusInternalServerError)
			return
		}
		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer clientConn.Close()

		if !isInternal {
			downR := http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Close:      true,
				ProtoMajor: r.ProtoMajor,
				ProtoMinor: r.ProtoMinor,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
				Header:     make(http.Header),
			}

			log.Println("We are connected, writing response")
			// if we don't close the body, write will just hang
			downR.Body.Close()
			downR.Write(clientConn)
		}

		data := netutils.HTTPReadWriteCloser{
			Writer:     pipW,
			ReadCloser: resp.Body,
		}

		netutils.RunPiper(clientConn, data)
		color.Yellow("Piper Finished")

		return
	}

	resp.ProtoMajor = r.ProtoMajor
	resp.ProtoMinor = r.ProtoMinor
	resp.Proto = r.Proto

	netutils.CopyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
}
