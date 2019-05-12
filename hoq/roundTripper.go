package hoq

import (
	"crypto/tls"
)

/**
快递员：负责将上层传入的请求发出并捕获对应的服务端响应
*/
type Courier interface {
	RoundTrip(*Request) (*Response, *RemoteInfo, error)
}

type TCPCourier struct {
}

func (t *TCPCourier) RoundTrip(req *Request) (resp *Response, remote *RemoteInfo, err error) {
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
	remote = &RemoteInfo{addr: conn.RemoteAddr()}
	resp, err = readResponse(conn)
	return
}

type QUICCourier struct {
	pool *quicConnPool
}

func (c *QUICCourier) RoundTrip(req *Request) (resp *Response, remote *RemoteInfo, err error) {
	if !req.ready() {
		return nil, nil, RequestNotReadyErr
	}
	s, err := c.pool.get(req.url.Host)
	sess := s.conn
	if err != nil {
		return
	}
	stream, err := sess.OpenStreamSync()
	if err != nil {
		return
	}
	defer stream.Close()
	err = req.Write(stream)
	if err != nil {
		return
	}
	remote = &RemoteInfo{addr: sess.RemoteAddr()}
	resp, err = readResponse(stream)
	return
}

func newQUICCourier() *QUICCourier {
	return &QUICCourier{
		pool: &quicConnPool{},
	}
}

func (c *QUICCourier) GetSession() {

}
