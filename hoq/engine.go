package hoq

import (
	"HOQ/logs"
	"errors"
	"github.com/lucas-clemente/quic-go"
	"net"
)

const (
	EngineTcp = iota
	EngineQuic
)

var UnsupportedEngine = errors.New("unsupported enginee")

/*
new a transporter according to it's name
*/
func newEngine(engine int, handler Handler) (engine, error) {
	switch engine {
	case EngineTcp:
		return newTcpEngine(handler), nil
	case EngineQuic:
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
	Name() string
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
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Error(err)
			continue
		}
		go t.HandleConnection(conn)
	}
}

func (t *tcpEngine) HandleConnection(con net.Conn) {
	logs.Debug("HandlerConnection in...")
	defer con.Close()
	req, err := readRequest(con)
	var resp *Response
	//error handler
	switch e := err.(type) {
	case nil:
		resp = t.handler(&Context{
			Request: req,
			Remote: &remoteInfo{
				addr: con.RemoteAddr(),
			},
		})
	case *ErrWithCode:
		resp, err = NewResponse(e.Code(), nil, nil)
		if err != nil {
			resp = innerServiceError
		}
	default:
		resp = innerServiceError
	}
	err = resp.Write(con)
	if err != nil {
		logs.Warn(err.Error())
		return
	}
	return
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
	var resp *Response
	//error handler
	switch e := err.(type) {
	case nil:
		resp = t.handler(&Context{
			Request: req,
			Remote: &remoteInfo{
				addr: sess.RemoteAddr(),
			},
		})
	case *ErrWithCode:
		resp, err = NewResponse(e.Code(), nil, nil)
		if err != nil {
			resp = innerServiceError
		}
	default:
		resp = innerServiceError
	}
	err = resp.Write(stream)
	if err != nil {
		logs.Warn(err.Error())
		return
	}
	return
}

func (t *quicEngine) Name() string {
	return "QUIC"
}

func (t *tcpEngine) Name() string {
	return "TCP"
}
