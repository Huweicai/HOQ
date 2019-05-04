package hoq

import (
	"bytes"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"io"
	"log"
)

/**
包装部分默认行为，简化API
*/
var defaultClient = Client{&QUICCourier{}}

/**
HTTP客户端，用于发起请求
*/o
type Client struct {
	engine Courier
}

/**
new common request
*/
func (c *Client) Request(targetUrl, method string, headers *Headers, body io.Reader) (ctx *Context, err error) {
	if !isSupportedMethod(method) {
		return nil, MethodNotSupportErr
	}
	u, err := urlParse(targetUrl)
	if err != nil {
		return
	}
	req := &Request{
		method:  method,
		Body:    body,
		headers: headers,
		proto:   defaultProto,
		url:     u,
	}
	resp, remoteInfo, err := c.engine.RoundTrip(req)
	ctx = &Context{
		Request:  req,
		Response: resp,
		Remote:   remoteInfo,
	}
	return
}

func (c *Client) Get(url string) (ctx *Context, err error) {
	return c.Request(url, MethodGET, nil, nil)
}

func (c *Client) Post(url string, body []byte) (ctx *Context, err error) {
	return c.Request(url, MethodPOST, nil, bytes.NewReader(body))
}

func (c *Client) Ping() bool {
	sess, err := quic.DialAddr(testHost, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	stream, err := sess.OpenStreamSync()
	if err != nil {
		log.Fatalf(err.Error())
		return false
	}
	_, err = stream.Write(testText)
	if err != nil {
		log.Fatalf(err.Error())
		return false
	}
	stream.Close()
	return true
}

func Get(url string) (ctx *Context, err error) {
	return defaultClient.Request(url, MethodGET, nil, nil)
}

func Post(url string, body []byte) (ctx *Context, err error) {
	return defaultClient.Request(url, MethodPOST, nil, bytes.NewReader(body))
}

func NewClient(engine int) (c *Client, err error) {
	switch engine {
	case TcpEngine:
		return &Client{&TCPCourier{}}, nil
	case QuicEngine:
		return &Client{&QUICCourier{}}, nil
	default:
		return nil, UnsupportedEngine
	}
}
