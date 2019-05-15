package hoq

import (
	"HOQ/logs"
	"bufio"
	"github.com/lucas-clemente/quic-go"
	_ "github.com/lucas-clemente/quic-go"
	"log"
)

type Server struct {
	engine  Engine
	addr    string
	handler Handler
}

/**
handler wraper
*/
func (s *Server) handle(ctx *Context) (resp *Response) {
	//pre
	resp = s.handler(ctx)
	//post
	resp.headers.GenDate()
	resp.headers.Set("X-Powered-By", "Go1.12.1")
	return
}

func handleRequestDemo(stream quic.Stream) {
	nr := bufio.NewReader(stream)
	got, _, err := nr.ReadLine()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	log.Println("GOT: " + string(got))
	_, err = stream.Write(testText)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	log.Println("SEND: " + string(testText))
}

/**
new a HTTP server with required arguments
*/
func NewServer(engine NGType, handler Handler) (s *Server, err error) {
	s = new(Server)
	ng, err := newEngine(engine, s)
	s.engine = ng
	s.handler = handler
	if err != nil {
		return
	}
	if handler == nil {
		return nil, ServerNotReadyErr
	}
	return
}

/**
check whether the config of server are all right for
the following running
*/
func (s *Server) Ready() bool {
	if s.engine == nil {
		return false
	}
	return true
}

func (s *Server) Run(addr string) error {
	if !s.Ready() {
		return ServerNotReadyErr
	}
	s.addr = addr
	logs.Info("server starting at", addr, "with engine", s.engine.Name())
	return s.engine.Serve(addr)
}

/**
start handle the request
*/
func work(channel Channel) {
	log.Println("start work with request from :" + channel.RemoteAddr().String())
	/*
		todo implements it
		先不用考虑连接复用的事情
		解析第一行；解析头部；封装成request，调用对应的handle方法
		返回response；flush；close
		finish
	*/
	return
}
