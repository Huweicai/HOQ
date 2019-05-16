package main

import (
	"HOQ/hoq"
	"time"
)

const urlQuic = "127.0.0.1:9002"
const urlTcp = "127.0.0.1:9001"

func main() {
	go startTwoServer()
	time.Sleep(time.Second)
	qc, tc := newClients()
	testClientSend(qc, urlQuic)
	testClientSend(tc, urlTcp)

}

func testClientSend(c *hoq.Client, u string) {
	ctx, err := c.Get("http://" + u)
	if err != nil {
		panic(err)
	}
	if ctx.Response.Code() != 200 {
		panic("not 200")
	}
}

func newClients() (*hoq.Client, *hoq.Client) {
	qc, err := hoq.NewClient(hoq.EngineQuic)
	if err != nil {
		panic(err.Error())
	}
	tc, err := hoq.NewClient(hoq.EngineTcp)
	if err != nil {
		panic(err.Error())
	}
	return qc, tc
}

func startTwoServer() {
	qserver, err := hoq.NewServer(hoq.EngineQuic, hoq.EchoHandler)
	if err != nil {
		panic(err)
	}
	tserver, err := hoq.NewServer(hoq.EngineTcp, hoq.EchoHandler)
	if err != nil {
		panic(err)
	}
	go qserver.Run(urlQuic)
	go tserver.Run(urlTcp)
}
