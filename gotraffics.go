package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	_          = iota
	KB float64 = 1 << (10 * iota)
	MB
)

var target string
var localPort string

func initFlag() {
	flag.StringVar(&target, "target", "", "target address")
	flag.StringVar(&localPort, "port", "8080", "local port to listen")
}

func handleRequestAndRedirect(resp http.ResponseWriter, req *http.Request) {
	bodySize := calculateSize(req.ContentLength)
	log.Printf("Url: %s ,\tbody size is: %s\n", req.URL.String(), bodySize)
	serveReverseProxy(resp, req)
}

func calculateSize(size int64) string {
	s := float64(size)
	switch {
	case s >= KB:
		return fmt.Sprintf("%.2fKB", s/KB)
	case s >= MB:
		return fmt.Sprintf("%.2fMB", s/MB)
	default:
		return fmt.Sprintf("%.2fB", s)
	}
}

// Serve a reverse proxy for a given url
func serveReverseProxy(resp http.ResponseWriter, req *http.Request) {
	// parse the url
	url, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Parse url error: %v", err)
	}

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host
	proxy.ServeHTTP(resp, req)
}

func logStartup() {
	log.Printf("Start proxy")
	log.Printf("Listen address is %v", "127.0.0.1:"+localPort)
	log.Printf("Target address is %v", target)
}

func main() {
	initFlag()
	flag.Parse()
	if flag.NFlag() == 0 {
		log.SetFlags(0)
		log.Fatal("Usage: gotraffics -target <target-address> -lport <listen-port>")
	}
	logStartup()
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(":"+localPort, nil); err != nil {
		panic(err)
	}
}
