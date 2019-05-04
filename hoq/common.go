package hoq

import (
	"io"
	"net/url"
)

var headerBodySepBytes = []byte("\r\n")

const headerBodySepStr = "\r\n"

const (
	defaultProto = "HTTP/1.1"

	ProtoHTTP  = "http"
	ProtoHTTPS = "https"
)

var NoBody = noBody{}

type noBody struct{}

func (noBody) Read([]byte) (int, error)         { return 0, io.EOF }
func (noBody) Close() error                     { return nil }
func (noBody) WriteTo(io.Writer) (int64, error) { return 0, nil }

/**
通用HTTP URL校验
*/
func urlParse(s string) (u *url.URL, err error) {
	u, err = url.Parse(s)
	if err != nil {
		return
	}
	//todo check Go HTTP库是不是也只支持这两个协议
	if u.Scheme != ProtoHTTP && u.Scheme != ProtoHTTPS {
		return nil, MalformedURLErr
	}
	if u.Host == "" {
		return nil, MalformedURLErr
	}
	return
}
