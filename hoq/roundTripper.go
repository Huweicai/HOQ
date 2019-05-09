package hoq

import (
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
)

/**
快递员：负责将上层传入的请求发出并捕获对应的服务端响应
*/
type Courier interface {
	RoundTrip(*Request) (*Response, *remoteInfo, error)
}

type TCPCourier struct {
}

func (t *TCPCourier) RoundTrip(req *Request) (resp *Response, remote *remoteInfo, err error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", req.url.Host, conf)
	if err != nil {
		return
	}
	err = req.Write(conn)
	if err != nil {
		return
	}
	remote = &remoteInfo{addr: conn.RemoteAddr()}
	resp, err = readResponse(conn)
	return
}

type QUICCourier struct {
}

func (c *QUICCourier) RoundTrip(req *Request) (resp *Response, remote *remoteInfo, err error) {
	if !req.ready() {
		return nil, nil, RequestNotReadyErr
	}
	sess, err := quic.DialAddr(req.url.Host, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return
	}
	stream, err := sess.OpenStreamSync()
	if err != nil {
		return
	}
	err = req.Write(stream)
	if err != nil {
		return
	}
	remote = &remoteInfo{addr: sess.RemoteAddr()}
	resp, err = readResponse(stream)
	return
}

func (c *QUICCourier) GetSession() {

}
