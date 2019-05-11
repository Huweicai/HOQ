package main

import (
	"HOQ/hoq"
)

func main() {
	qserver, err := hoq.NewServer(hoq.EngineQuic, hoq.EchoHandler)
	if err != nil {
		panic(err)
	}
	tserver, err := hoq.NewServer(hoq.EngineTcp, hoq.EchoHandler)
	if err != nil {
		panic(err)
	}
	go qserver.Run(":6665")
	go tserver.Run(":6667")
	c := make(chan interface{})
	<-c
}
