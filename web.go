package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
}

func run() error {
	allow, err := regexp.Compile(os.Getenv("ALLOW"))
	if err != nil {
		return err
	}

	p, err := newProxy(allow, os.Getenv("REMOTE_URL"))
	if err != nil {
		return err
	}

	if err := http.ListenAndServe(":3000", p); err != nil {
		return err
	}

	return nil
}

type proxy struct {
	allow  *regexp.Regexp
	proxy  http.Handler
	remote string
}

func newProxy(allow *regexp.Regexp, remote string) (*proxy, error) {
	proxy := &proxy{allow: allow, remote: remote}

	u, err := url.Parse(remote)
	if err != nil {
		return nil, err
	}

	rp := httputil.NewSingleHostReverseProxy(u)

	rp.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	proxy.proxy = rp

	return proxy, nil
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s -> %s%s: ", r.RemoteAddr, p.remote, r.URL.Path)

	if !p.allow.MatchString(r.URL.Path) {
		http.Error(w, "not allowed", 403)
		fmt.Println("deny")
		return
	}

	fmt.Println("allow")

	p.proxy.ServeHTTP(w, r)
}
