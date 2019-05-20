package hoq

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

/**
Request 是单纯的应用层，HTTP维度的Request
只承载HTTP维度的信息
todo new request func
*/
type Request struct {
	method  string
	url     *url.URL
	proto   string
	headers *Headers
	Body    io.Reader
}

func (r *Request) FirstLine() string {
	return r.requestLine()
}

func (r *Request) URL() *url.URL {
	return r.url
}

func (r *Request) Method() string {
	return r.method
}

func (r *Request) EatFirstLine(s string) error {
	method, rawUrl, proto, ok := parseFirstRequestLine(s)
	if !ok {
		return NewErrWithCode(StatusBadRequest, "malformated HTTP request")
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		return WrapErrWithCode(StatusBadRequest, err)
	}
	r.method = method
	r.url = u
	r.proto = proto
	return nil
}

func (r *Request) GetBody() io.Reader {
	return r.Body
}

func (r *Request) SetBody(body io.Reader) {
	r.Body = body
}

func (r *Request) SetHeader(h *Headers) {
	r.headers = h
}

func (r *Request) GetHeader() *Headers {
	return r.headers
}

func NewRequest(method, targetUrl string, headers *Headers, body io.Reader) (r *Request, err error) {
	if !IsSupportedMethod(method) {
		return nil, MethodNotSupportErr
	}
	u, err := urlParse(targetUrl)
	if err != nil {
		return
	}
	if headers == nil {
		headers = NewHeaders(nil)
	}
	headers.GenContentLength(body)
	headers.GenCookies(u)
	r = &Request{
		method:  method,
		Body:    body,
		headers: headers,
		proto:   defaultProto,
		url:     u,
	}
	return
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
func readRequest(reader io.Reader) (r *Request, err error) {
	r = &Request{}
	err = read(reader, r)
	if err != nil {
		return
	}
	err = validateRequest(r)
	return
}

/**
根据RFC定义，校验request格式
*/
func validateRequest(r *Request) error {
	//method post 未标注ContentLength，
	// (Section 3.3.2 of [RFC7230])
	if r.method == MethodPOST && !r.headers.Exits(HeaderContentLength) {
		return NewErrWithCode(StatusLengthRequired, "header length is required")
	}
	//RFC  6.5.11.  413 Payload Too Large
	if r.headers.ContentLength() > 300000 {
		return NewErrWithCode(StatusRequestEntityTooLarge, "Payload Too Large")
	}
	//RFC 6.5.15.  426 Upgrade Required
	//HTTP/1.0 HTTP/0.9
	if r.proto < minProto {
		return NewErrWithCode(StatusUpgradeRequired, "Upgrade Required")
	}
	if r.proto > maxProto {
		return NewErrWithCode(StatusHTTPVersionNotSupported, "version too high too support")
	}
	return nil
}

/**
read body into bytes not just get the body in Reader format
*/
func (r *Request) ReadBody() (body []byte, err error) {
	if !existBody(r.Body) {
		return nil, nil
	}
	reader := r.Body
	switch r.headers.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
	case "flate":
		reader = flate.NewReader(r.Body)
	case "zlib":
		reader, err = zlib.NewReader(r.Body)
	default:
	}
	if err != nil {
		return
	}
	return ioutil.ReadAll(reader)
}

func (r *Request) Response(code int, headers *Headers, body []byte) (rsp *Response, err error) {
	return NewResponse(code, headers, bytes.NewReader(body))
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
	_, err = io.WriteString(writer, r.headers.Serialize()+headerBodySepStr+headerBodySepStr)
	if err != nil {
		return
	}
	if !existBody(r.Body) {
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

func (r *Request) wrap() (*http.Request, error) {
	return http.NewRequest(r.method, r.url.String(), r.Body)
}
