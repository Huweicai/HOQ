package hoq

import (
	"HOQ/logs"
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
			logs.Error(err.Error())
			continue
		}
		go t.HandleSess(sess)
	}
}

/**
一个Quic session 对应一条底层UDP连接，
但是一个Session上理论可以被很多HTTP连接复用
*/
func (t *quicEngine) HandleSess(sess quic.Session) {
	logs.Debug("HandleSess in")
	stream, err := sess.AcceptStream()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	defer sess.Close()
	t.HandleStream(sess, stream)
}

/**
todo 错误请求对应状态码返回
*/
func (t *quicEngine) HandleStream(sess quic.Session, stream quic.Stream) {
	logs.Debug("HandleStream in")
	defer stream.Close()
	req, err := readRequest(stream)
	if err != nil {
		logs.Warn("read request failed", err)
		return
	}
	resp := t.handler(&Context{
		Request: req,
		Remote: &remoteInfo{
			addr: sess.RemoteAddr(),
		},
	})
	err = resp.Write(stream)
	if err != nil {
		logs.Warn(err.Error())
		return
	}
	return
}
