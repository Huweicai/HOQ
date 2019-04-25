package hoq

import (
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
)

type Courier interface {
	RoundTrip(*Request) (*Response, error)
}

type TCPCourier struct {
}

func (t *TCPCourier) RoundTrip(*Request) (*Response, error) {
	panic("implements me")
	return nil, nil
}

type QUICCourier struct {
}

func (c *QUICCourier) RoundTrip(req *Request) (resp *Response, err error) {
	txt, err := req.Serialize()
	if err != nil {
		return
	}
	sess, err := quic.DialAddr(req.url.Host, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return
	}
	stream, err := sess.OpenStreamSync()
	if err != nil {
		return
	}
	_, err = stream.Write([]byte(txt))
	if err != nil {
		return
	}
	return newResponse(stream)
}

func (c *QUICCourier) GetSession() {

}
