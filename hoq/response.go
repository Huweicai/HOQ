package hoq

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	headers    *Headers
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
		headers:    &Headers{mimeHeaderToMap(mimeHeader)},
	}, nil
}

/**
序列化成文本便于传输
todo body是否需要单独抽出来？
或许序列化先只需要序列化报文头部分
Deprecated
*/
func (r *Response) Serialize() (b []byte, err error) {
	headerLine := r.statusLine()
	headers := r.headers.Serialize()
	b = []byte(headerLine + headerBodySepStr + headers + headerBodySepStr + headerBodySepStr)
	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return b, err
		}
		b = append(b, body...)
	}
	return b, nil
}

/*
check whether response is ready for transport
*/
func (r *Response) ready() bool {
	if r.proto == "" {
		return false
	}
	if 0 < r.statusCode || r.statusCode > 600 {
		return false
	}
	if r.statusMSg == "" {
		return false
	}
	return true
}

/**
状态行
*/
func (r *Response) statusLine() string {
	code := strconv.Itoa(r.statusCode)
	if code == "0" {
		code = ""
	}
	return fmt.Sprintf("%s %s %s", r.proto, code, r.statusMSg)
}

func (r *Response) Write(writer io.Writer) (err error) {
	if !r.ready() {
		return ResponseNotReadyErr
	}
	//1st
	_, err = io.WriteString(writer, r.statusLine()+headerBodySepStr)
	if err != nil {
		return
	}
	//2nd
	hds := r.headers.Serialize()
	_, err = io.WriteString(writer, hds+headerBodySepStr)
	if err != nil {
		return
	}
	//3rd
	if r.Body == nil {
		return
	}
	_, err = io.Copy(writer, r.Body)
	if err != nil {
		return
	}
	return
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
