package hoq

import (
	"bufio"
	"github.com/lucas-clemente/quic-go"
	_ "github.com/lucas-clemente/quic-go"
	"log"
)

type Server struct {
	engine  engine
	addr    string
	handler Handler
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
	listen, err := s.engine.Listen(addr)
	if err != nil {
		return err
	}
	log.Println("HTTP server started ,listen at " + addr)
	for {
		channel, err := listen.Accept()
		if err != nil {
			return err
		}
		go work(channel)
	}
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
