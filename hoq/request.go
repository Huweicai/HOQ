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
	headers *Headers
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
	//todo how to handler the situation without Header
	mimeHeader, err := textR.ReadMIMEHeader()
	if err != nil {
		return
	}
	//todo judge body exists

	return &Request{
		method:  method,
		headers: &Headers{headers: mimeHeaderToMap(mimeHeader)},
		url:     u,
		proto:   proto,
		Body:    stream,
	}, nil
}

func (r *Request) GetBody() ([]byte, error) {
	return ioutil.ReadAll(r.Body)
}

func (r *Request) Response(code int, headers *Headers, body []byte) *Response {
	if headers == nil {
		headers = &Headers{}
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
Deprecated
*/
func (r *Request) Serialize() (b []byte, err error) {
	headerLine := fmt.Sprintf("%s %s %s", r.method, r.url.String(), r.proto)
	b = []byte(headerLine + headerBodySepStr)
	if r.headers != nil {
		headers := r.headers.Serialize()
		b = append(b, []byte(headers+headerBodySepStr+headerBodySepStr)...)
	}
	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		b = append(b, body...)
	}
	return
}

/**
请求行
RFC 7230           HTTP/1.1 Message Syntax and Routing         June 2014
 A request-line begins with a method token, followed by a single space
   (SP), the request-target, another single space (SP), the protocol
   version, and ends with CRLF.

     request-line   = method SP request-target SP HTTP-version CRLF

   The method token indicates the request method to be performed on the
   target resource.  The request method is case-sensitive.

     method         = token

   The request methods defined by this specification can be found in
   Section 4 of [RFC7231], along with information regarding the HTTP
   method registry and considerations for defining new methods.

   The request-target identifies the target resource upon which to apply
   the request, as defined in Section 5.3.

   Recipients typically parse the request-line into its component parts
   by splitting on whitespace (see Section 3.5), since no whitespace is
   allowed in the three components.  Unfortunately, some user agents
   fail to properly encode or exclude whitespace found in hypertext
   references, resulting in those disallowed characters being sent in a
   request-target.

   Recipients of an invalid request-line SHOULD respond with either a
   400 (Bad Request) error or a 301 (Moved Permanently) redirect with
   the request-target properly encoded.  A recipient SHOULD NOT attempt
   to autocorrect and then process the request without a redirect, since
   the invalid request-line might be deliberately crafted to bypass
   security filters along the request chain.

   HTTP does not place a predefined limit on the length of a
   request-line, as described in Section 2.5.  A server that receives a
   method longer than any that it implements SHOULD respond with a 501
   (Not Implemented) status code.  A server that receives a
   request-target longer than any URI it wishes to parse MUST respond
   with a 414 (URI Too Long) status code (see Section 6.5.12 of
   [RFC7231]).
   Various ad hoc limitations on request-line length are found in
   practice.  It is RECOMMENDED that all HTTP senders and recipients
   support, at a minimum, request-line lengths of 8000 octets.
*/
func (r *Request) requestLine() string {
	u := ""
	if r.url != nil {
		u = r.url.String()
	}
	return fmt.Sprintf("%s %s %s", r.method, u, r.proto)
}

/**
将请求序列化写入对应的Writer中
*/
func (r *Request) Write(writer io.Writer) (err error) {
	//第一步：request line
	_, err = io.WriteString(writer, r.requestLine()+headerBodySepStr)
	if err != nil {
		return
	}
	//第二步：header
	//header todo 分重要级，需要把Host ， Connection , UA , Content-Length 等控制字段放置于最前方
	_, err = io.WriteString(writer, r.headers.Serialize()+headerBodySepStr)
	if err != nil {
		return
	}
	if r.Body == nil {
		return
	}
	//第三步：body（如果有的话）
	_, err = io.Copy(writer, r.Body)
	if err != nil {
		return
	}
	return
}

/*
todo implements it
*/
func (r *Request) ready() bool {
	return true
}
