package hoq

//refer to RFC 7321 https://tools.ietf.org/html/rfc7231#page-24
const (
	MethodGET     = "GET"
	MethodHead    = "HEAD"
	MethodPOST    = "POST"
	MethodPUT     = "PUT"
	MethodDELETE  = "DELETE"
	MethodCONNECT = "CONNECT"
	MethodOPTIONS = "OPTIONS"
	MethodTRACE   = "TRACE"
)

var methods = []string{MethodGET, MethodHead, MethodPOST, MethodPUT, MethodDELETE, MethodCONNECT, MethodOPTIONS, MethodOPTIONS, MethodTRACE}

/**
是否在支持的八种方法中
*/
func isSupportedMethod(method string) bool {
	for _, me := range methods {
		if me == method {
			return true
		}
	}
	return false
}
