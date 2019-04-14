package hoq

import (
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
	Body    []byte
}

//parse the first line of header
//such as "GET /foo HTTP/1.1"
func parseFirstLine(line string) (method, url, proto string, ok bool) {
	i1 := strings.Index(line, " ")
	i2 := strings.Index(line[i1+1:], " ")
	if i1 == -1 || i2 == -1 {
		return
	}
	i2 = i1 + i2 + 1
	return line[:i1], line[i1+1 : i2], line[i2+1:], true
}
