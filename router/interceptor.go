package router

import (
	"HOQ/hoq"
	"HOQ/logs"
	"fmt"
	"strconv"
)

type PreInterceptor func(ctx *hoq.Context, r *Router) error
type SufInterceptor func(resp *hoq.Response, r *Router) error

/**
日志记录请求信息
*/
func LogPreInterceptor(ctx *hoq.Context, r *Router) error {
	logs.Info("request in ", blue(ctx.Request.Method()), ctx.Request.URL().Path)
	return nil
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(96), s)
}

const (
	RecordPrefixMethod   = "m_"
	RecordPrefixPath     = "p_"
	RecordPrefixRespCode = "r_"
)

func RecordPreInterceptor(ctx *hoq.Context, r *Router) error {
	method := ctx.Request.Method()
	path := ctx.Request.URL().Path
	r.recordChan <- &record{RecordPrefixMethod + method, 1}
	r.recordChan <- &record{RecordPrefixPath + path, 1}
	return nil
}

func RecordSufInterceptor(rsp *hoq.Response, r *Router) error {
	code := strconv.Itoa(rsp.Code())
	r.recordChan <- &record{RecordPrefixRespCode + code, 1}
	return nil
}

type record struct {
	k string
	i int
}
