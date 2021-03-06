package hoq

import (
	"io"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//todo 限制头部字段单个字段大小不得超过9000，请求中最多有120个字段
const (
	HeaderContentLength = "Content-Length"
	HeaderHost          = "Host"
	HeaderUserAgent     = "User-Agent"
	HeaderConnection    = "Connection"
	HeaderDate          = "Date"
)

/**
头部字段，不应该仅仅只是一个map
针对Host等字段需要单独定义
*/
type Headers struct {
	headers    map[string]string
	gotCookies []*Cookie
	setCookies []*Cookie
}

/**
生成Date字段
*/
func (h *Headers) GenDate() {
	now := time.Now().UTC().String()
	h.headers[HeaderDate] = now
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

func (h *Headers) GenCookies(u *url.URL) {
	cks := jar.GetCookies(u)
	h.gotCookies = append(h.gotCookies, cks...)
}

/**
从Reader中读入Header
*/
func ReadHeaders(reader *textproto.Reader) (*Headers, error) {
	mimeHeader, err := reader.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	setCookies := readSetCookies(mimeHeader)
	gotCookies := readGotCookies(mimeHeader)
	m := mimeHeaderToMap(mimeHeader)
	return &Headers{headers: m, setCookies: setCookies, gotCookies: gotCookies}, nil
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

/**
判断Header是否存在
*/
func (h *Headers) Exits(k string) bool {
	_, ok := h.headers[k]
	return ok
}
func (h *Headers) Serialize() string {
	b := strings.Builder{}
	b.WriteString(h.serializeGotCookies())
	b.WriteString(h.serializeSetCookies())
	//todo check omit values after the first one
	for name, value := range h.headers {
		if value == "" {
			continue
		}
		b.WriteString(name + ": " + value + headerBodySepStr)
	}
	return strings.TrimSuffix(b.String(), headerBodySepStr)
}

func (h *Headers) GenContentLength(body io.Reader) bool {
	l, ok := bodyLength(body)
	if !ok {
		return false
	}
	h.headers[HeaderContentLength] = strconv.Itoa(int(l))
	return true
}

func (h *Headers) Set(k, v string) *Headers {
	h.headers[k] = v
	return h
}

func (h *Headers) Get(k string) string {
	got, _ := h.headers[k]
	return got
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

func convertHttpHeader(hd http.Header) *Headers {
	m := make(map[string]string)
	for k, vs := range hd {
		if len(vs) == 0 {
			continue
		}
		m[k] = vs[0]
	}
	return NewHeaders(m)
}

func (h *Headers) serializeSetCookies() string {
	out := ""
	for _, ck := range h.gotCookies {
		out += "Set-Cookie: " + ck.String() + headerBodySepStr
	}
	return out
}

func (h *Headers) serializeGotCookies() string {
	if len(h.gotCookies) == 0 {
		return ""
	}
	out := "Cookies: "
	for _, ck := range h.gotCookies {
		out += ck.KVString() + headerBodySepStr
	}
	return out
}

func (h *Headers) AddSetCookie(c *Cookie) {
	h.setCookies = append(h.setCookies, c)

}

func (h *Headers) AddGotCookie(c *Cookie) {
	h.gotCookies = append(h.gotCookies, c)
}
