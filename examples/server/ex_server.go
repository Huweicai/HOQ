package main

import (
	"HOQ/hoq"
	"log"
)

func main() {
	server, _ := hoq.NewServer(hoq.EngineQuic, TestHandler)
	err := server.Run("127.0.0.1:8787")
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func TestHandler(ctx *hoq.Context) *hoq.Response {
	rsp, _ := ctx.Request.Response(hoq.StatusOK, nil, []byte("Bye Bye~"))
	rsp.AddHeader("Set-Cookie", "abc=123;domain=127.0.0.1")
	return rsp
}
