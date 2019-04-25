package hoq

type Handler func(ctx *Context) *Response

/**
a simple echo handler for test
*/
func EchoHandler(ctx *Context) *Response {
	body, _ := ctx.Request.GetBody()
	if len(body) == 0 {
		body = []byte("HELLO WORLD")
	}
	return ctx.Request.Response(StatusOK, nil, body)
}
