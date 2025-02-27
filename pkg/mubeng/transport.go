package mubeng

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
        "h12.io/socks"
)

// Transport to auto-switch transport between HTTP/S or SOCKSv5 proxies.
// Depending on the protocol scheme, returning value of http.Transport with Dialer or Proxy.
func Transport(p string) (tr *http.Transport, err error) {
	proxyURL, err := url.Parse(p)
	if err != nil {
		return nil, err
	}

	switch proxyURL.Scheme {
	case "socks4":
          	dialer := socks.Dial("socks4://"+proxyURL.Host)
		tr = &http.Transport{
			Dial: dialer,
		}
	case "socks5":
	        dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, getAuth(proxyURL), proxy.Direct)
	        if err != nil {
	     	     return nil, err
	        }

		tr = &http.Transport{
			Dial: dialer.Dial,
		}
	case "http":
		tr = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	default:
		return nil, fmt.Errorf("unsupported proxy protocol scheme: %s", proxyURL.Scheme)
	}

	tr.DisableKeepAlives = true
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return tr, nil
}

func getAuth(URL *url.URL) *proxy.Auth {
	auth := &proxy.Auth{}
	user := URL.User.Username()
	pass, _ := URL.User.Password()

	if user != "" {
		auth.User = user

		if pass != "" {
			auth.Password = pass
		}
	}

	return auth
}
