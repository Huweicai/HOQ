package main

import (
	"HOQ/hoq"
	"HOQ/logs"
)

var headers = map[string]string{"User-Agent": "Firefox"}

func main() {
	nc, _ := hoq.NewClient(hoq.QuicEngine)

	ctx, err := nc.Request(hoq.MethodGET, "http://127.0.0.1:8787", hoq.NewHeaders(headers), nil)
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
