package hoq

import (
	"bytes"
	"io"
	"net/url"
	"strings"
)

var headerBodySepBytes = []byte("\r\n")

const headerBodySepStr = "\r\n"

const (
	minProto     = "HTTP/1.0"
	maxProto     = "HTTP/2.0"
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
probe whether body is exist
*/
func existBody(r io.Reader) bool {
	return r != nil && r != NoBody
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

/**
probe length for reader
*/
func bodyLength(r io.Reader) (length int64, ok bool) {
	switch i := r.(type) {
	case nil, noBody:
		return 0, true
	case *bytes.Reader:
		return int64(i.Len()), true
	case *strings.Reader:
		return int64(i.Len()), true
	case *bytes.Buffer:
		return int64(i.Len()), true
	default:
		return
	}
}

// bodyAllowedForStatus reports whether a given response status code
// permits a body. See RFC 7230, section 3.3.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == 204:
		return false
	case status == 304:
		return false
	}
	return true
}
