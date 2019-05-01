package hoq

import (
	"log"
	"net/url"
)

const (
	headerBodySep = byte(10)
	defaultProto  = "HTTP/1.1"

	ProtoHTTP  = "http"
	ProtoHTTPS = "https"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

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
