package hoq

import "errors"

var (
	ServerNotReadyErr   = errors.New("some configs may not right for running")
	MethodNotSupportErr = errors.New("method not support")
	MalformedURLErr     = errors.New("malformed url")
	ResponseNotReadyErr = errors.New("response not valid ")
	RequestNotReadyErr  = errors.New("request not valid")
)
