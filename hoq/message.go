package hoq

import "io"

//一种报文，分为Request 和 Response
type Message interface {
	FirstLine() string
	EatFirstLine(s string) error

	GetBody() io.Reader
	SetBody(body io.Reader)

	SetHeader(h *Headers)
	GetHeader() *Headers
}
