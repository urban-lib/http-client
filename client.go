package httpClient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client interface {
	ProxyType() string
	ProxyHost() string
	Request(method, uri string, data interface{}, header http.Header, resp chan Response)
}

type client struct {
	host  string
	proxy struct {
		t    string
		host string
	}
}

func NewClient(host string, proxyType, proxyHost string) Client {
	return &client{
		host: host,
		proxy: struct {
			t    string
			host string
		}{t: proxyType, host: proxyHost},
	}
}

func (c *client) ProxyType() string {
	return c.proxy.t
}

func (c *client) ProxyHost() string {
	return c.proxy.host
}

func (c *client) determineProxy() (*http.Client, error) {
	switch c.proxy.t {
	case "socks5":
		return NewSocks5Client(c.proxy.host)
	default:
		return nil, fmt.Errorf("Proxy type [%s] not supported", c.proxy.t)
	}
}

func (c *client) getClient() (*http.Client, error) {
	if c.proxy.host != "" {
		return c.determineProxy()
	}
	return &http.Client{}, nil
}

func (c *client) Request(method, uri string, data interface{}, header http.Header, resp chan Response) {
	switch method {
	case http.MethodPost:
		c.post(uri, data.([]uint8), header, resp)
		return
	default:
		resp <- Response{
			Status: http.StatusMethodNotAllowed,
			Body:   nil,
			Error:  fmt.Errorf("Method [%s] not allowed!!", method),
		}
		return
	}
}

func (c *client) post(uri string, body []byte, header http.Header, response chan Response) {
	cli, err := c.getClient()
	if err != nil {
		response <- Response{
			Status: http.StatusTeapot,
			Body:   nil,
			Error:  err,
		}
		return
	}
	request, err := http.NewRequest(http.MethodPost, c.host+uri, bytes.NewBuffer(body))
	if err != nil {
		response <- Response{
			Status: http.StatusTeapot,
			Body:   nil,
			Error:  err,
		}
		return
	}
	request.Header = header
	resp, err := cli.Do(request)
	if err != nil {
		response <- Response{
			Status: http.StatusTeapot,
			Body:   nil,
			Error:  err,
		}
		return
	}
	defer resp.Body.Close()
	bodyResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response <- Response{
			Status: http.StatusTeapot,
			Body:   nil,
			Error:  err,
		}
		return
	}
	response <- Response{
		Status: resp.StatusCode,
		Body:   bodyResponse,
		Error:  err,
	}
	return
}
