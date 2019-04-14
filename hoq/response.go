package hoq

import "net/url"

/**
Request 是单纯的应用层，HTTP维度的Request
只承载Response维度的信息
*/
type Response struct {
	method  string
	url     *url.URL
	proto   string
	headers Headers
	Body    []byte
	request *Request
}
