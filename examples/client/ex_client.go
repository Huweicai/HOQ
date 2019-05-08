package main

import (
	"HOQ/hoq"
	"HOQ/logs"
)

var headers = map[string]string{"User-Agent": "Firefox"}

func main() {
	nc, _ := hoq.NewClient(hoq.EngineQuic)
	//ctx, err := nc.Request(hoq.MethodGET, "http://127.0.0.1:8787", hoq.NewHeaders(headers), strings.NewReader("666666666"))
	ctx, err := nc.Get("http://127.0.0.1:8787/bye")
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info(ctx.Response.FirstLine())
	got, err := ctx.Response.ReadBody()
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info(string(got))
}
