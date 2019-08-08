package main

import (
	"HOQ/hoq"
	"HOQ/logs"
	"HOQ/router"
	"strconv"
	"time"
)

const urlServerQuic = "127.0.0.1:9001"
const urlServerTcp = "127.0.0.1:9002"
const urlRouterQuic = "127.0.0.1:9003"
const urlRouterTcp = "127.0.0.1:9004"

func main() {
	testClientAndRouter()
}

func testClientAndRouter() {
	startTwoRouter()
	qc, tc := newClients()
	testRouterClient(qc, urlRouterQuic)
	testRouterClient(tc, urlRouterTcp)
}

func testClientAndServer() {
	//server 和 client测试
	startTwoServer()
	time.Sleep(time.Second)
	qc, tc := newClients()
	testClientSend(qc, urlServerQuic)
	testClientSend(tc, urlServerTcp)
}

func testRouterClient(c *hoq.Client, u string) {
	u = "http://" + u
	ctx, err := c.Get(u + "/200")
	if err != nil {
		panic(err)
	}
	if ctx.Response.Code() != hoq.StatusOK {
		panic("not 200")
	}
	ctx, err = c.Post(u+"/200", nil)
	if ctx.Response.Code() != hoq.StatusMethodNotAllowed {
		panic("not 200" + strconv.Itoa(ctx.Response.Code()))
	}
	ctx, err = c.Request(hoq.MethodGET, u+"/allowAllMethods", nil, nil)
	if err != nil {
		panic(err)
	}
	if ctx.Response.Code() != hoq.StatusOK {
		panic("not 200")
	}
	ctx, err = c.Get(u + "/noutFound")
	if err != nil {
		panic(err)
	}
	if ctx.Response.Code() != hoq.StatusNotFound {
		panic("can not be found")
	}
	ctx, err = c.Request(hoq.MethodPUT, u+"/allowAllMethods", nil, nil)
	if ctx.Response.Code() != hoq.StatusOK {
		panic("not 200")
	}
	//超时
	c.SetReqTimeout(100 * time.Millisecond)
	ctx, err = c.Get(u + "/timeout")
	if err != hoq.RequestTimeoutErr {
		panic("should timeout")
	}
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
	go qserver.Run(urlServerQuic)
	go tserver.Run(urlServerTcp)
}

func startTwoRouter() {
	qr, err := router.New(hoq.EngineQuic, urlRouterQuic, 0)
	if err != nil {
		panic(err)
	}
	qr.GET("/200", hoq.EchoHandler)
	qr.Add("/allowAllMethods", hoq.EchoHandler, hoq.MethodsWildCard)
	qr.Add("/timeout", sleepHandler, hoq.MethodsWildCard)
	tr, err := router.New(hoq.EngineTcp, urlRouterTcp, 0)
	if err != nil {
		panic(err)
	}
	tr.GET("/200", hoq.EchoHandler)
	tr.Add("/allowAllMethods", hoq.EchoHandler, hoq.MethodsWildCard)
	tr.Add("/timeout", sleepHandler, hoq.MethodsWildCard)
	go qr.Run()
	go tr.Run()
}

func sleepHandler(ctx *hoq.Context) *hoq.Response {
	body, _ := ctx.Request.ReadBody()
	if len(body) == 0 {
		body = []byte("HELLO WORLD")
	}
	time.Sleep(500 * time.Millisecond)
	logs.Info(string(body))
	rsp, _ := ctx.Request.Response(hoq.StatusOK, nil, body)
	return rsp
}
