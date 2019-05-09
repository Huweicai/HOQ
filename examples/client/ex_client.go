package main

import (
	"HOQ/hoq"
	"HOQ/logs"
	"crypto/tls"
	"net/http"
)

var headers = map[string]string{"User-Agent": "Firefox"}

func main() {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	rsp, err := c.Get("http://127.0.0.1:8787")
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info(rsp)

	nc, _ := hoq.NewClient(hoq.EngineTcp)
	//ctx, err := nc.Request(hoq.MethodGET, "http://127.0.0.1:8787", hoq.NewHeaders(headers), strings.NewReader("666666666"))
	for i := 0; i < 1000; i++ {
		ctx, err := nc.Head("https://127.0.0.1:8787/hello")
		ctx, err = nc.Get("https://127.0.0.1:8787/bye")
		ctx, err = nc.Post("https://127.0.0.1:8787/hello", nil)
		if err != nil {
			logs.Error(err)
			continue
		}
		logs.Info(ctx.Response.FirstLine())
		got, err := ctx.Response.ReadBody()
		if err != nil {
			logs.Error(err)
			return
		}
		logs.Info(string(got))
	}
}
