package main

import (
	"HOQ/hoq"
	"HOQ/logs"
)

func main() {
	nc, _ := hoq.NewClient(hoq.QuicEngine)
	ctx, err := nc.Request("http://127.0.0.1:8787", hoq.MethodGET, &hoq.Headers{
		map[string][]string{
			"": []string{},
		},
	}, nil)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Info(ctx.Response)
}
