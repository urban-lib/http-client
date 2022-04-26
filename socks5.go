package httpClient

import (
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"runtime"
	"time"
)

type DialContext func(ctx context.Context, network, address string) (net.Conn, error)

func NewSocks5Client(host string) (*http.Client, error) {
	baseDialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	var dialContext DialContext
	if host != "" {
		dialSocksProxy, err := proxy.SOCKS5("tcp", host, nil, baseDialer)
		if err != nil {
			return nil, fmt.Errorf("Error creating SOCKS5 proxy. %v", err)
		}
		if contextDieler, ok := dialSocksProxy.(proxy.ContextDialer); ok {
			dialContext = contextDieler.DialContext
		} else {
			return nil, fmt.Errorf("Failed type assertion to DialContext")
		}
	} else {
		dialContext = (baseDialer).DialContext
	}

	return proxyClient(dialContext, host), nil
}

func proxyClient(dialContext DialContext, host string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 nil,
			DialContext:           dialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0),
		},
	}
}
