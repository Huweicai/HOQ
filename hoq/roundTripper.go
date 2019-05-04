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

func (t *TCPCourier) RoundTrip(*Request) (*Response, *remoteInfo, error) {
	panic("implements me")
	return nil, nil, nil
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
	resp, err = newResponse(stream)
	return
}

func (c *QUICCourier) GetSession() {

}
