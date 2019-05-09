package hoq

import (
	"bufio"
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

func (r *Response) Code() int {
	return r.statusCode
}

func (r *Response) FirstLine() string {
	return r.statusLine()
}

func (r *Response) EatFirstLine(s string) error {
	var ok bool
	r.statusCode, r.statusMSg, r.proto, ok = parseStatusLine(s)
	if !ok {
		return NewErrWithCode(StatusBadRequest, "malformed HTTP response")
	}
	return nil
}

func (r *Response) GetBody() io.Reader {
	return r.Body
}

func (r *Response) SetBody(body io.Reader) {
	r.Body = body
}

func (r *Response) SetHeader(h *Headers) {
	r.headers = h
}

func (r *Response) GetHeader() *Headers {
	return r.headers
}

/**
常用Response
*/
var innerServiceError *Response = FastResponse(StatusInternalServerError, nil)

/**
不会返回Error的简单初始化Response api
*/
func FastResponse(code int, headers *Headers) *Response {
	if headers == nil {
		headers = NewHeaders(nil)
	}
	msg := StatusMessage(code)
	return &Response{
		proto:      defaultProto,
		statusCode: code,
		statusMSg:  msg,
		headers:    headers,
		Body:       NoBody,
	}
}

/**
新建一个Response并初始化
*/
func NewResponse(code int, headers *Headers, body io.Reader) (rsp *Response, err error) {
	if headers == nil {
		headers = NewHeaders(nil)
	}
	msg := StatusMessage(code)
	headers.GenContentLength(body)
	return &Response{
		proto:      defaultProto,
		statusCode: code,
		statusMSg:  msg,
		headers:    headers,
		Body:       body,
	}, nil
}

/**
core function!!!
*/
func readResponse(reader io.Reader) (r *Response, err error) {
	r = &Response{}
	err = read(reader, r)
	return
}

func read(reader io.Reader, r Message) (err error) {
	bufR := bufio.NewReader(reader)
	textR := textproto.NewReader(bufR)
	line, err := textR.ReadLine()
	if err != nil {
		return
	}
	//1st
	err = r.EatFirstLine(line)
	if err != nil {
		return
	}
	//2nd
	headers, err := ReadHeaders(textR)
	if err != nil {
		return
	}
	//3rd
	var body io.Reader = NoBody
	if length := headers.ContentLength(); length > 0 {
		body = io.LimitReader(bufR, length)
	}
	//return
	r.SetHeader(headers)
	r.SetBody(body)
	return
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
	if 0 > r.statusCode || r.statusCode > 600 {
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
	_, err = io.WriteString(writer, hds+headerBodySepStr+headerBodySepStr)
	if err != nil {
		return
	}
	//3rd
	if !existBody(r.Body) {
		return
	}
	//_, err = io.WriteString(writer, "HELLO WORLD")
	_, err = io.Copy(writer, r.Body)
	if err != nil {
		return
	}
	return
}

//parse the first line of response header
//such as "HTTP/1.1 200 OK"
func parseStatusLine(line string) (code int, msg, proto string, ok bool) {
	i1 := strings.Index(line, " ")
	if i1 < 0 {
		ok = false
		return
	}
	proto = line[:i1]
	tmp := strings.Index(line[i1+1:], " ")
	if i1 < -1 {
		ok = false
		return
	}
	i2 := tmp + i1 + 1
	codeS := line[i1+1 : i2]
	code, err := strconv.Atoi(codeS)
	if err != nil {
		return
	}
	msg = line[i2+1:]
	ok = true
	return
}

func (r *Response) ReadBody() ([]byte, error) {
	if !existBody(r.Body) {
		return nil, nil
	}
	return ioutil.ReadAll(r.Body)
}
