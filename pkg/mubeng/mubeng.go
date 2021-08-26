package mubeng

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"log"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// New define HTTP client & request of http.Request itself.
//
// also removes Hop-by-hop headers when it is sent to backend (see http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html),
// then add X-Forwarded-For header value with the IP address value of rotator proxy IP.
func stripCtlAndExtFromUnicode(str string) string {
	isOk := func(r rune) bool {
		return r < 32 || r >= 127
	}
	// The isOk filter is such that there is no need to chain to norm.NFC
	t := transform.Chain(norm.NFKD, transform.RemoveFunc(isOk))
	// This Transformer could also trivially be applied as an io.Reader
	// or io.Writer filter to automatically do such filtering when reading
	// or writing data anywhere.
	str, _, _ = transform.String(t, str)
	return str
}

func (proxy *Proxy) New(req *http.Request) (*http.Client, *http.Request) {
	client = &http.Client{Transport: proxy.Transport}

	// http: Request.RequestURI can't be set in client requests.
	// http://golang.org/src/pkg/net/http/client.go
	req.RequestURI = ""

	for _, h := range HopHeaders {
		req.Header.Del(h)
	}
	proxyURL, err := url.Parse(stripCtlAndExtFromUnicode(proxy.Address))	
	log.Printf("Proxy: %s",proxy.Address)
	if err != nil {
		//log.Fatal(err)
		log.Printf("%s",err)
	}
	if host, _, err := net.SplitHostPort(proxyURL.Host); err == nil {
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			host = strings.Join(prior, ", ") + ", " + host
		}
		req.Header.Set("X-Forwarded-For", host)
	}

	req.Header.Set("X-Forwarded-Proto", req.URL.Scheme)
	return client, req
}
