package hoq

import "time"

type Context struct {
	RemoteAddr string
	Request    *Request
	Response   *Response
}

func (Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (Context) Done() <-chan struct{} {
	return nil
}

func (Context) Err() error {
	return nil
}

func (Context) Value(key interface{}) interface{} {
	return ""
}
