package hoq

import (
	"bufio"
	"errors"
	"github.com/lucas-clemente/quic-go"
	"github.com/sirupsen/logrus"
)

const (
	TcpEngine = iota
	QuicEngine
)

var UnsupportedEngine = errors.New("unsupported enginee")

/*
new a transporter according to it's name
*/
func newEngine(engine int, handler Handler) (engine, error) {
	switch engine {
	case TcpEngine:
		return newTcpEngine(handler), nil
	case QuicEngine:
		return newQuicEngine(handler), nil
	default:
		return nil, UnsupportedEngine
	}
}

/**
  底层运输载体，目前支持tcp , quic 两种
*/
type engine interface {
	Serve(addr string) error
}

type tcpEngine struct {
	handler Handler
}

func newTcpEngine(handler Handler) *tcpEngine {
	return &tcpEngine{
		handler: handler,
	}
}

type quicEngine struct {
	handler Handler
}

func newQuicEngine(handler Handler) *quicEngine {
	return &quicEngine{
		handler: handler,
	}
}

func (t *tcpEngine) Serve(addr string) error {
	panic("implement me")
}

func (t *quicEngine) Serve(addr string) error {
	listen, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	for {
		sess, err := listen.Accept()
		if err != nil {
			logrus.Error(err.Error())
			continue
		}
		go t.HandleSess(sess)
	}
}
func (t *quicEngine) HandleSess(sess quic.Session) {
	stream, err := sess.AcceptStream()
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	defer sess.Close()
	t.HandleStream(sess, stream)
}

func (t *quicEngine) HandleStream(sess quic.Session, stream quic.Stream) {
	defer stream.Close()
	req, err := readRequest(stream)
	resp := t.handler(&Context{
		Request:    req,
		RemoteAddr: sess.RemoteAddr().String(),
	})
	rspText, err := resp.Serialize()
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	_, err = bufio.NewWriter(stream).Write(rspText)
	if err != nil {
		logrus.Warn(err.Error())
		return
	}
	return
}
