package hoq

import (
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"log"
)

type Client struct {
}

func (*Client) Get(url string) *Response {
	panic("no implemented")
}

func (*Client) Post(url string, body []byte) *Response {
	panic("no implemented")
}

func (c *Client) Ping() bool {
	sess, err := quic.DialAddr(testHost, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		log.Fatalf(err.Error())
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
