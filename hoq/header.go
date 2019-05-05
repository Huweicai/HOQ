package hoq

import (
	"io"
	"net/textproto"
	"strconv"
)

const (
	HeaderContentLength = "Content-Length"
	HeaderHost          = "Host"
	HeaderUserAgent     = "User-Agent"
	HeaderConnection    = "Connection"
)

/**
头部字段，不应该仅仅只是一个map
针对Host等字段需要单独定义
*/
type Headers struct {
	headers map[string]string
}

/**
新建Header from map
*/
func NewHeaders(m map[string]string) *Headers {
	if m == nil {
		m = make(map[string]string)
	}
	return &Headers{headers: m}
}

/**
从Reader中读入Header
*/
func ReadHeaders(reader *textproto.Reader) (*Headers, error) {
	mimeHeader, err := reader.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	m := mimeHeaderToMap(mimeHeader)
	return &Headers{m}, nil
}

/**
返回Header中申明的Content-Length字段
在ReadRequest阶段，Content-Length就需要校验，不合法的CL 需要返回对应的bad code
*/
func (h *Headers) ContentLength() int64 {
	l, ok := h.headers[HeaderContentLength]
	//没设置即视为不存在Body
	if !ok {
		return 0
	}
	i, err := strconv.Atoi(l)
	if err != nil {
		return 0
	}
	if i < 0 {
		return 0
	}
	return int64(i)
}

func (h *Headers) Serialize() string {
	out := ""
	//todo check omit values after the first one
	for name, value := range h.headers {
		if value == "" {
			continue
		}
		out += name + ": " + value
	}
	return out
}

func (h *Headers) GenContentLength(body io.Reader) bool {
	l, ok := bodyLength(body)
	if !ok {
		return false
	}
	h.headers[HeaderContentLength] = strconv.Itoa(int(l))
	return true
}

func (h *Headers) Set(k, v string) {

}

func (h *Headers) Get(k string) string {
	//todo implements it
	return ""
}

func mimeHeaderToMap(mime textproto.MIMEHeader) map[string]string {
	if len(mime) == 0 {
		return nil
	}
	m := make(map[string]string)
	for key, values := range mime {
		if len(values) == 0 {
			continue
		}
		m[key] = values[0]
	}
	return m
}
