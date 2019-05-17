package main

import (
	"HOQ/hoq"
	"HOQ/logs"
	"time"
)

var headers = map[string]string{"User-Agent": "Firefox"}

func main() {
	nc, _ := hoq.NewClient(hoq.EngineQuic)
	nc.SetReqTimeout(1 * time.Second)
	//ctx, err := nc.Request(hoq.MethodGET, "http://127.0.0.1:8787", hoq.NewHeaders(headers), strings.NewReader("666666666"))
	for i := 0; i < 1000; i++ {
		ctx, err := nc.Get("https://127.0.0.1:8787/bye")
		if err != nil {
			logs.Error(err)
		}
		got, err := ctx.Response.ReadBody()
		if err != nil {
			logs.Error(err)
		}
		logs.Info(string(got))
	}
}
