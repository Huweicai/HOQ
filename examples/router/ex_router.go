package main

import (
	"HOQ/hoq"
	"HOQ/logs"
	"HOQ/router"
	"time"
)

func main() {
	r, err := router.New(hoq.EngineTcp, "127.0.0.1:8787", 0)
	if err != nil {
		logs.Error(err)
		return
	}
	r.Add("/hello", hoq.EchoHandler, hoq.MethodGET)
	r.Add("/bye", hoq.ByeHandler, hoq.MethodGET)
	go r.ShowRecordsCyclic(1 * time.Second)
	logs.Error(r.Run())
}
