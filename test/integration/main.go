package main

import (
	"HOQ/hoq"
	"HOQ/logs"
	"time"
)

const urlQuic = "127.0.0.1:9002"
const urlTcp = "127.0.0.1:9001"

func main() {
	//server 和 client测试
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
	//测试正常情况
	if ctx.Response.Code() != 200 {
		panic("not 200")
	}
	logs.Info("正常GET请求测试通过")
	txt := "<HTML><H1>HELLO WORLD</H1></HTML>"
	ctx, err = c.Post("http://"+u, []byte(txt))
	if err != nil {
		panic(err.Error())
	}
	logs.Info("正常POST请求测试通过")
	got, err := ctx.Response.ReadBody()
	if err != nil {
		panic(err.Error())
	}
	if string(got) != txt {
		panic("body not right")
	}
	logs.Info("Body处理解析成功")
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
