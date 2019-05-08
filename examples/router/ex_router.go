package main

import (
	"HOQ/hoq"
	"HOQ/logs"
	"HOQ/router"
)

func main() {
	r, err := router.New(hoq.EngineQuic, "127.0.0.1:8787", 0)
	if err != nil {
		logs.Error(err)
		return
	}
	r.Add("/hello", hoq.EchoHandler, hoq.MethodGET)
	r.Add("/bye", hoq.ByeHandler, hoq.MethodGET)
	logs.Error(r.Run())
}
