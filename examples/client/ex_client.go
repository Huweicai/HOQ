package main

import (
	"HOQ/hoq"
	"HOQ/logs"
)

var headers = map[string]string{"User-Agent": "Firefox"}

func main() {
	nc, _ := hoq.NewClient(hoq.EngineTcp)
	//ctx, err := nc.Request(hoq.MethodGET, "http://127.0.0.1:8787", hoq.NewHeaders(headers), strings.NewReader("666666666"))
	ctx, err := nc.Get("http://127.0.0.1:8787")
	if err != nil {
		logs.Error(err)
		return
	}
	got, err := ctx.Response.ReadBody()
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info(string(got))
}
