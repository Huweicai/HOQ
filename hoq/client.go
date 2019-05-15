package hoq

import (
	"HOQ/logs"
	"bytes"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"io"
	"log"
	"time"
)

const defaultReqTimeout = 5 * time.Second

/**
包装部分默认行为，简化API
*/
var defaultClient = Client{engine: &QUICCourier{}, reqTimeout: defaultReqTimeout}

/**
HTTP客户端，用于发起请求
*/
type Client struct {
	//超时时长
	reqTimeout time.Duration
	engine     Courier
}

/**
设置请求超时时间
*/
func (c *Client) SetReqTimeout(t time.Duration) {
	c.reqTimeout = t
}

/**
new common request
*/
func (c *Client) Request(method, targetUrl string, headers *Headers, body io.Reader) (ctx *Context, err error) {
	req, err := NewRequest(method, targetUrl, headers, body)
	if err != nil {
		return
	}
	req.headers.GenDate()
	doneChan := make(chan int)
	var resp *Response
	var remoteInfo *RemoteInfo
	go func() {
		resp, remoteInfo, err = c.engine.RoundTrip(req)
		doneChan <- 0
	}()
	select {
	case <-doneChan:
		logs.Debug("request finished before timeout")
	case <-time.After(c.reqTimeout):
		err = RequestTimeoutErr
		return
	}
	ctx = &Context{
		Request:  req,
		Response: resp,
		Remote:   remoteInfo,
	}
	return
}

func (c *Client) Get(url string) (ctx *Context, err error) {
	return c.Request(MethodGET, url, nil, nil)
}

func (c *Client) Head(url string) (ctx *Context, err error) {
	return c.Request(MethodHead, url, nil, nil)
}

func (c *Client) Post(url string, body []byte) (ctx *Context, err error) {
	return c.Request(MethodPOST, url, nil, bytes.NewReader(body))
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

func NewClient(engine NGType) (c *Client, err error) {
	switch engine {
	case EngineTcp:
		return &Client{engine: &TCPCourier{}, reqTimeout: defaultReqTimeout}, nil
	case EngineQuic:
		return &Client{engine: newQUICCourier(), reqTimeout: defaultReqTimeout}, nil
	default:
		return nil, UnsupportedEngine
	}
}
