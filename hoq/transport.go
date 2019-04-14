package hoq

import (
	"errors"
	"github.com/lucas-clemente/quic-go"
)

const (
	TcpEngine = iota
	QuicEngine
)

var UnsupportedEngine = errors.New("unsupported enginee")

/*
new a transporter according to it's name
*/
func newTransporter(engine int) (engine, error) {
	switch engine {
	case TcpEngine:
		return new(tcpEngine), nil
	case QuicEngine:
		return new(quicEngine), nil
	default:
		return nil, UnsupportedEngine
	}
}

/**
  底层运输载体，目前支持tcp , quic 两种
*/
type engine interface {
	Listen(addr string) (httpListener, error)
	Run() error
}
type tcpEngine struct {
}
type quicEngine struct {
}

func (t *tcpEngine) Listen(addr string) (httpListener, error) {
	panic("implement me")
}

func (t *tcpEngine) Run() error {
	panic("implement me")
}

func (quicEngine) Listen(addr string) (httpListener, error) {
	lsl, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return nil, err
	}
	return
}

func (q *quicEngine) Run() error {
	q.Listen()
}
