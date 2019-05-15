package hoq

import "errors"

var (
	ServerNotReadyErr   = errors.New("some configs may not right for running")
	MethodNotSupportErr = errors.New("method not support")
	MalformedURLErr     = errors.New("malformed url")
	ResponseNotReadyErr = errors.New("response not valid ")
	RequestNotReadyErr  = errors.New("request not valid")

	ConnectTimeoutErr = errors.New("connection timeout")
	RequestTimeoutErr = errors.New("request timeout")
)

type ErrWithCode struct {
	code int
	msg  string
}

func WrapErrWithCode(code int, err error) *ErrWithCode {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	return &ErrWithCode{code, msg}
}

func NewErrWithCode(code int, msg string) *ErrWithCode {
	return &ErrWithCode{code, msg}
}

func (e *ErrWithCode) Error() string {
	return e.msg
}

func (e *ErrWithCode) Code() int {
	return e.code
}
