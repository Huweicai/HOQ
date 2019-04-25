package hoq

import (
	"bufio"
	"errors"
	"io"
	"net/textproto"
	"strconv"
	"strings"
)

/**
Request 是单纯的应用层，HTTP维度的Request
只承载Response维度的信息
example:
HTTP/1.1 200 OK
Server: nginx
Date: Thu, 25 Apr 2019 11:46:43 GMT
Content-Type: text/plain;charset=UTF-8
Cache-Control: no-store

this is body
*/
type Response struct {
	proto      string
	statusCode int
	statusMSg  string
	headers    Headers
	Body       io.Reader
}

func newResponse(reader io.Reader) (r *Response, err error) {
	bufR := bufio.NewReader(reader)
	textR := textproto.NewReader(bufR)
	line, err := textR.ReadLine()
	if err != nil {
		return
	}
	code, msg, proto, ok := parseFirstResponseLine(line)
	if !ok {
		return nil, errors.New("malformed HTTP response")
	}
	mimeHeader, err := textR.ReadMIMEHeader()
	if err != nil {
		return
	}
	return &Response{
		proto:      proto,
		statusCode: code,
		statusMSg:  msg,
		Body:       reader,
		headers:    Headers(mimeHeader),
	}, nil
}

func (r *Response) Serialize() (b string, err error) {
	//todo implements it
	return `HTTP/1.1 200 OK
Server: nginx
Date: Thu, 25 Apr 2019 11:46:43 GMT
Content-Type: text/plain;charset=UTF-8
Cache-Control: no-store

this is body`, nil
}

//parse the first line of response header
//such as "HTTP/1.1 200 OK"
func parseFirstResponseLine(line string) (code int, msg, proto string, ok bool) {
	i1 := strings.Index(line, " ")
	proto = line[:i1]
	i2 := strings.Index(line[i1+1:], " ")
	code, err := strconv.Atoi(line[i1+1 : i2])
	if err != nil {
		return
	}
	msg = line[i2+1:]
	ok = true
	return
}
