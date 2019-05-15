package hoq

type Handler func(ctx *Context) *Response

/**
a simple echo handler for test
*/
func EchoHandler(ctx *Context) *Response {
	body, _ := ctx.Request.ReadBody()
	if len(body) == 0 {
		body = []byte("HELLO WORLD")
	}
	rsp, _ := ctx.Request.Response(StatusOK, nil, body)
	return rsp
}

func ByeHandler(ctx *Context) *Response {
	rsp, _ := ctx.Request.Response(StatusOK, nil, []byte("Bye Bye~"))
	return rsp
}
