package hoq

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/textproto"
	"net/url"
	"strings"
)

/**
Request 是单纯的应用层，HTTP维度的Request
只承载HTTP维度的信息
*/
type Request struct {
	method  string
	url     *url.URL
	proto   string
	headers Headers
	Body    io.Reader
}

//parse the first line of header
//such as "GET /foo HTTP/1.1"
func parseFirstRequestLine(line string) (method, url, proto string, ok bool) {
	i1 := strings.Index(line, " ")
	i2 := strings.Index(line[i1+1:], " ")
	if i1 == -1 || i2 == -1 {
		return
	}
	i2 = i1 + i2 + 1
	return line[:i1], line[i1+1 : i2], line[i2+1:], true
}

/**
convert a io.Reader to a HTTP request
*/
func readRequest(stream io.Reader) (r *Request, err error) {
	bufR := bufio.NewReader(stream)
	textR := textproto.NewReader(bufR)
	fl, err := textR.ReadLine()
	if err != nil {
		return
	}
	method, rawUrl, proto, ok := parseFirstRequestLine(fl)
	u, err := url.Parse(rawUrl)
	if err != nil {
		return
	}
	if !ok {
		return nil, errors.New("malformed HTTP request")
	}
	mimeHeader, err := textR.ReadMIMEHeader()
	if err != nil {
		return
	}

	return &Request{
		method:  method,
		headers: Headers(mimeHeader),
		url:     u,
		proto:   proto,
		Body:    stream,
	}, nil
}

func (r *Request) GetBody() ([]byte, error) {
	return ioutil.ReadAll(r.Body)
}

func (r *Request) Response(code int, headers Headers, body []byte) *Response {
	if headers == nil {
		headers = make(Headers)
	}
	msg := StatusMessage(code)
	return &Response{
		proto:      r.proto,
		statusCode: code,
		statusMSg:  msg,
		headers:    headers,
		Body:       bytes.NewReader(body),
	}
}

/*
序列化
*/
func (r *Request) Serialize() (b []byte, err error) {
	headerLine := fmt.Sprintf("%s %s %s", r.method, r.url.String(), r.proto)
	b = []byte(headerLine + "\n")
	if r.headers != nil {
		headers := r.headers.Serialize()
		b = append(b, []byte(headers+"\n")...)
	}
	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		b = append(b, headerBodySep)
		b = append(b, body...)
	}
	return
}
